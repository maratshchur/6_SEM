from caesar_cipher import caesar_cipher
from vigenere_cipher import vigenere_cipher

def main():
    filename = input("Введите название файла (например, input.txt): ")
    
    try:
        with open(filename, "r", encoding="utf-8") as file:
            input_text = file.read()
    except FileNotFoundError:
        print(f"Файл '{filename}' не найден.")
        return
    except UnicodeDecodeError:
        print(f"Ошибка декодирования файла '{filename}'. Проверьте кодировку файла.")
        return

    action_choice = input("Выберите действие (1 - Зашифровать, 2 - Расшифровать): ")

    if action_choice == "1":
        cipher_choice = input("Выберите метод шифрования (1 - Цезарь, 2 - Виженер): ")

        if cipher_choice == "1":
            key = int(input("Введите ключ для шифра Цезаря (целое число): "))
            ciphered_text = caesar_cipher(input_text, key)
            ciphered_filename = "caesar_ciphered.txt"
        elif cipher_choice == "2":
            key = input("Введите ключ для шифра Виженера (строка): ")
            ciphered_text = vigenere_cipher(input_text, key)
            ciphered_filename = "vigenere_ciphered.txt"
        else:
            print("Некорректный выбор метода шифрования.")
            return

        with open(ciphered_filename, "w", encoding="utf-8") as file:
            file.write(ciphered_text)
        print(f"Зашифрованный текст записан в файл: {ciphered_filename}")

    elif action_choice == "2":
        cipher_choice = input("Выберите метод расшифрования (1 - Цезарь, 2 - Виженер): ")

        if cipher_choice == "1":
            key = int(input("Введите ключ для шифра Цезаря (целое число): "))
            deciphered_text = caesar_cipher(input_text, key, encrypt=False)
            deciphered_filename = "caesar_deciphered.txt"
        elif cipher_choice == "2":
            key = input("Введите ключ для шифра Виженера (строка): ")
            deciphered_text = vigenere_cipher(input_text, key, encrypt=False)
            deciphered_filename = "vigenere_deciphered.txt"
        else:
            print("Некорректный выбор метода расшифрования.")
            return

        with open(deciphered_filename, "w", encoding="utf-8") as file:
            file.write(deciphered_text)
        print(f"Расшифрованный текст записан в файл: {deciphered_filename}")

    else:
        print("Некорректный выбор действия.")

if __name__ == "__main__":
    main()