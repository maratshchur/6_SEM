#!/bin/bash

board=(. . . . . . . . .)

display_board() {
  clear
  echo " ${board[0]} | ${board[1]} | ${board[2]} "
  echo "---+---+---"
  echo " ${board[3]} | ${board[4]} | ${board[5]} "
  echo "---+---+---"
  echo " ${board[6]} | ${board[7]} | ${board[8]} "
}

screenshot() {
  timestamp=$(date +%Y%m%d_%H%M%S)
  filename="screenshot_${timestamp}.txt"
  echo " ${board[0]} | ${board[1]} | ${board[2]} " > "$filename"
  echo "---+---+---" >> "$filename"
  echo " ${board[3]} | ${board[4]} | ${board[5]} " >> "$filename"
  echo "---+---+---" >> "$filename"
  echo " ${board[6]} | ${board[7]} | ${board[8]} " >> "$filename"
  echo "Скриншот сохранен в файл: $filename"
}

check_winner() {
  local player=$1
  local winning_combinations=( "0 1 2" "3 4 5" "6 7 8" "0 3 6" "1 4 7" "2 5 8" "0 4 8" "2 4 6" )

  for combination in "${winning_combinations[@]}"; do
    local -a indices=($combination)
    if [[ ${board[${indices[0]}]} == "$player" && ${board[${indices[1]}]} == "$player" && ${board[${indices[2]}]} == "$player" ]]; then
      return 0
    fi
  done

  return 1
}

reset_game() {
  board=(. . . . . . . . .)
  player=X
  moves=0
}

reset_game
while true; do
  display_board
  echo "Ход игрока $player. Введите номер ячейки (1-9) или 's' для скриншота:"
  
  read -r move
  
  if [[ "$move" == "s" ]]; then
    screenshot
    continue
  fi

  if [[ ! "$move" =~ ^[1-9]$ ]]; then
    echo "Неверный ввод. Пожалуйста, введите число от 1 до 9."
    sleep 1
    continue
  fi

  cell=$((move - 1))

  if [[ "${board[cell]}" != "." ]]; then
    echo "Ячейка занята. Пожалуйста, выберите другую."
    sleep 1  
    continue
  fi

  board[cell]=$player
  moves=$((moves + 1))

  if check_winner "$player"; then
    display_board
    echo "Поздравляем! Игрок $player победил!"
    echo "Хотите сыграть снова? (y/n)"
    read -r play_again
    if [[ "$play_again" == "y" ]]; then
      reset_game
      continue
    else
      echo "Спасибо за игру!"
      break
    fi
  elif [[ $moves -eq 9 ]]; then
    display_board
    echo "Ничья!"
    echo "Хотите сыграть снова? (y/n)"
    read -r play_again
    if [[ "$play_again" == "y" ]]; then
      reset_game
      continue
    else
      echo "Спасибо за игру!"
      break
    fi
  fi

  if [[ "$player" == "X" ]]; then
    player=O
  else
    player=X
  fi
done