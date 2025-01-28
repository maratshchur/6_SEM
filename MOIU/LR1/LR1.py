import numpy as np

def compute_modified_inverse(A_inv: np.ndarray, vec: np.ndarray, index: int) -> np.ndarray:
    num_rows, num_cols = A_inv.shape

    # Шаг 1: Умножаем обратную матрицу на столбец x
    modified_vector = A_inv @ vec
    if modified_vector[index] == 0:
        raise ValueError(f"Matrix is not invertible, since modified_vector[{index}] == 0")

    # Шаг 2: Создаем копию вектора l и заменяем i-й элемент на -1
    original_value = modified_vector[index]
    modified_vector_copy = modified_vector.copy()
    modified_vector_copy[index] = -1

    # Шаг 3: Умножаем каждое число получившегося вектора на -1/l[i]
    scaled_vector = (-1 / original_value) * modified_vector_copy

    # Шаг 4: Создаём еденичную матрицу с i-м столбцом равным l с шапочкой
    transformation_matrix = np.eye(num_rows)
    transformation_matrix[:, index] = scaled_vector

    # Шаг 5: Оптимизируем нахождение произведения двух матриц
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
    # Пример из условия
    matrix_a = np.array([[1, -1, 0],
                         [0, 1, 0],
                         [0, 0, 1]])
    inverse_a = np.linalg.inv(matrix_a)
    vector_x = np.array([1, 0, 1])
    index_i = 2
    print(compute_modified_inverse(inverse_a, vector_x, index_i))

    #//
    matrix_b = np.array([[1, 0, 5],
                         [2, 1, 6],
                         [3, 4, 0]])

    inverse_b = np.linalg.inv(matrix_b)
    vector_y = np.array([2, 2, 2])
    index_j = 1
    print(compute_modified_inverse(inverse_b, vector_y, index_j))

    # Необратимая матрица
    # matrix_c = np.array([[2, 1], [4, 1]])
    # inverse_c = np.linalg.inv(matrix_c)
    # vector_z = np.array([1,2])
    # index_k = 1
    # print(compute_modified_inverse(inverse_c, vector_z, index_k))

if __name__ == "__main__":
    main()