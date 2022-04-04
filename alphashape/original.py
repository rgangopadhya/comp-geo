from collections import defaultdict

import numpy as np
from scipy.spatial import Delaunay
from shapely.geometry import mapping, MultiLineString, MultiPoint, Point, Polygon
from shapely.ops import polygonize, unary_union


ALPHA = 300


def alphashape(points, alpha=None):
    # Adapted from https://git.io/JvXFR for our use case here
    """
    Compute the alpha shape (concave hull) of a set of points.  If the number
    of points in the input is three or less, the convex hull is returned to the
    user.  For two points, the convex hull collapses to a `LineString`; for one
    point, a `Point`.
    Args:
        points (list or ``shapely.geometry.MultiPoint`` or
            ``geopandas.GeoDataFrame``): an iterable container of points
        alpha (float): alpha value
    Returns:
        ``shapely.geometry.Polygon`` or ``shapely.geometry.LineString`` or
        ``shapely.geometry.Point`` or ``geopandas.GeoDataFrame``: the resulting
            geometry
    """
    # If given a triangle for input, or an alpha value of zero or less,
    # return the convex hull.
    if len(points) < 4 or alpha is None or alpha <= 0:
        points = MultiPoint(list(points))
        return points.convex_hull

    coords = np.array(points)
    tri = Delaunay(coords)
    inv_alpha = 1.0 / alpha

    # indices of points making up triangles in the triangulation (pylint disable is for known pylint
    # issue with Delaunay class)
    simplices = tri.simplices  # pylint: disable=no-member

    # Coordinates of each vertex of the triangles
    tri_coords = coords[simplices]

    # Lengths of sides of triangle
    a = np.sqrt(
        np.sum(np.square(tri_coords[:, 0, :] - tri_coords[:, 1, :]), 1))
    b = np.sqrt(
        np.sum(np.square(tri_coords[:, 1, :] - tri_coords[:, 2, :]), 1))
    c = np.sqrt(
        np.sum(np.square(tri_coords[:, 2, :] - tri_coords[:, 0, :]), 1))
    # Semiperimeter of triangle
    s = (a + b + c) * 0.5
    # Area of triangle by Heron's formula
    area = s * (s - a) * (s - b) * (s - c)
    circumradius = a * b * c / (4.0 * np.sqrt(area))
    # Filter based on the circumradius value compared with the alpha value:
    # 1. circumradius is less than inv_alpha -> triangle is included in shape
    # 2. circumradius is finite and >= inv_alpha -> triangle should not be included in shape
    # 3. circumradius is NaN or infinite -> ignore (Due to how we handle holes below, it's
    #    simpler to just skip these triangles altogether.)
    included_triangles = circumradius < inv_alpha
    hole_triangles = np.all(
        [np.isfinite(circumradius), circumradius >= inv_alpha], 0)

    # If none of the triangles are included in the shape, fallback to convex hull (rather than
    # returning empty geometry)
    if not np.any(included_triangles):
        points = MultiPoint(list(points))
        return points.convex_hull

    # Count how many faces (triangles) included in alphashape each edge has (0, 1, or 2)
    internal_faces_count = defaultdict(int)
    for ia, ib, ic in simplices[included_triangles]:
        for i, j in [(ia, ib), (ib, ic), (ic, ia)]:
            ordered_indices = (i, j) if i < j else (j, i)
            internal_faces_count[ordered_indices] += 1

    # Collect faces (triangles) that are not in alphashape
    holes = [Polygon([pa, pb, pc])
             for pa, pb, pc in tri_coords[hole_triangles]]

    # Compute the alphashape by finding all the edges on the boundary of the shape, which is all
    # edges which make up exactly one triangle included in the alphashape. However, since we don't
    # know whether each of these edges should be part of the exterior ring(s) of the shape or
    # interior rings (holes), we treat them all as exterior (by polygonizing + unioning them). We
    # then difference out the holes (those triangles excluded from the alphashape), to arrive at the
    # correct alphashape.
    #
    # NOTE: Conceptually, this is equivalent to computing the union of all triangles included in the
    # alphashape (and what the library this is copied from does). However, since there are often far
    # more triangles included in the shape than excluded, it is significantly faster doing it as
    # outlined above instead. This is also equivalent to computing the convex hull and differencing
    # out the holes from that, but the above is slightly faster
    #
    # TODO: We shouldn't need to do any unioning at all here - if we properly track the half edges
    # (oriented edges associated with a single face), we can track half edges on the boundary and
    # determine directly if the rings they form are CW or CCW winding to determine if they are a
    # hole or not. But doing this is fast enough for now, so not making this optimization yet.
    boundary_edges = []
    for (i, j), count in internal_faces_count.items():
        assert count in (1, 2), f"({i},{j}): {count}"
        if count == 1:
            boundary_edges.append(coords[[i, j]])

    m = MultiLineString(boundary_edges)
    exterior_polygons = unary_union(list(polygonize(m)))

    valid_holes = []
    for hole in holes:
        if not hole.is_valid:
            print(f"invalid hole: {mapping(hole)}")
            hole = hole.buffer(0)
        valid_holes.append(hole)

    try:
        hole_polygons = unary_union(valid_holes)
    except ValueError as e:
        holes_json = [mapping(hole) for hole in valid_holes]
        print(
            f"Error constructing holes multipolygon:\nholes: {holes_json}\nexception: {e}"
        )
        raise e

    return exterior_polygons.difference(hole_polygons)
