import numpy as np
import time

def compute_modified_inverse(A_inv: np.ndarray, vec: np.ndarray, index: int) -> np.ndarray:
    num_rows, num_cols = A_inv.shape

    modified_vector = A_inv @ vec
    if modified_vector[index] == 0:
        raise ValueError(f"Matrix is not invertible, since modified_vector[{index}] == 0")

    original_value = modified_vector[index]
    modified_vector_copy = modified_vector.copy()
    modified_vector_copy[index] = -1

    scaled_vector = (-1 / original_value) * modified_vector_copy

    transformation_matrix = np.eye(num_rows)
    transformation_matrix[:, index] = scaled_vector

    result_inverse = optimized_matrix_multiplication(transformation_matrix, A_inv, index)

    return result_inverse

def optimized_matrix_multiplication(transformation_matrix: np.ndarray, A_inv: np.ndarray, index: int) -> np.ndarray:
    num_rows, num_cols = A_inv.shape
    result = np.zeros((num_rows, num_cols))

    for row in range(num_rows):
        for col in range(num_cols):
            product = transformation_matrix[row][row] * A_inv[row][col]
            if row == index:
                result[row][col] = product
            else:
                additional_product = transformation_matrix[row][index] * A_inv[index][col]
                result[row][col] = product + additional_product

    return result

def main():
    size = 10000
    matrix_a = np.random.rand(size, size)
    
    matrix_a += size * np.eye(size)  

    start_time = time.time()
    inverse_a_numpy = np.linalg.inv(matrix_a)
    numpy_time = time.time() - start_time
    print(f"NumPy inversion time: {numpy_time:.4f} seconds")

    vector_x = np.random.rand(size)
    index_i = np.random.randint(0, size)

    start_time = time.time()
    inverse_a_custom = compute_modified_inverse(inverse_a_numpy, vector_x, index_i)
    custom_time = time.time() - start_time
    print(f"Custom method inversion time: {custom_time:.4f} seconds")

if __name__ == "__main__":
    main()