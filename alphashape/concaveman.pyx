# cython
# distutils: language = c++
import numpy as np


cdef extern from "concaveman.cpp":
    void pyconcaveman2d(double *points_c, size_t num_points, int * hull_points_c, size_t num_hull_points, double concavity, double lengthThreshold, double **concave_points_c, size_t *num_concave_points, void(**p_free)(void*))


cdef extern from "concaveman.h":
    pass


def concaveman2d(points, hull, concavity=2.0, lengthThreshold=0.0):
    points = np.array(points).astype(np.double)
    hull = np.array(hull).astype(np.int32)

    if len(points.shape) != 2:
        raise ValueError('points must be a 2-D array')

    if len(hull.shape) != 1:
        raise ValueError('hull must be a 1-D array')

    if np.any(hull >= len(points)) or np.any(hull < 0):
        raise ValueError('hull indices out of bounds')

    # p_concave_points_c = _ffi.new('double**')
    # p_num_concave_points = _ffi.new('size_t*')
    # p_free = _ffi.new('void(**)(void*)')
    cdef double** p_concave_points_c
    cdef size_t *p_num_concave_points
    cdef void(**p_free)(void*)

    # points_c = _ffi.cast('double*', points.ctypes.data)
    # hull_c = _ffi.cast('int*', hull.ctypes.data)
    pyconcaveman2d(
        points,
        len(points),
        hull,
        len(hull),
        concavity,
        lengthThreshold,
        p_concave_points_c,
        p_num_concave_points,
        p_free,
    )

    num_concave_points = p_num_concave_points[0]
    concave_points_c = p_concave_points_c[0]

    # buffer = _ffi.buffer(concave_points_c, 8 * 2 * num_concave_points)

    concave_points = np.frombuffer(concave_points_c, dtype=np.double)
    concave_points = concave_points.reshape((num_concave_points, 2))
    concave_points = concave_points.copy()

    print('concave_points:', concave_points)

    p_free[0](concave_points_c)

    return concave_points
