# coding: utf8
u"""Модуль консольного пользовательского интерфейса."""


import input_data
import collection_data
import processing


def main_menu():
    u"""Отображает основное меню программы."""
    def show_actions():
        u"""Выводит доступные действия."""
        actions = [
            "Start"
        ]
        for i, action in enumerate(actions):
            output = "COMPUTER [Main menu]: " + str(i + 1) + " == " + action
            print(output)
        print("COMPUTER [Main menu]: 00 == Quit")
        number_actions = len(actions)
        return number_actions
    
    def check_user_answer(number_actions):
        u"""Получение ответа пользователя."""
        output = "USER [Main menu]: (1-" + str(number_actions) + "/00) "
        user_answer = raw_input(output)
        if user_answer == "00":
            print("COMPUTER: Quit...")
            exit(0)
        elif user_answer == "1":
            run_processing()
        else:
            print("COMPUTER: Error. Pleace, repeat input.")
            return check_user_answer(number_actions)

    number_actions = show_actions()
    check_user_answer(number_actions)


def run_processing():
    u"""Запускает функцию обработки."""
    def show_begin_points():
        u"""Отображает точки начала сбора подписчиков."""
        begin_points = [
            "From old to recent (asc)", "From recent to old (desc)"
        ]
        for i, begin_point in enumerate(begin_points):
            output = "COMPUTER [.. -> Begin from]: " + str(i + 1) + " == " + \
                begin_point
            print(output)
    
    def check_user_answer():
        u"""Получение пользовательского ответа."""
        begin_point = input_data.get_sort()
        if begin_point == "1":
            begin_point = "time_asc"
            return True
        elif begin_point == "2":
            begin_point = "time_desc"
            return True
        else:
            print("COMPUTER: Error. Pleace, repeat input.")
            return check_user_answer()

    values = collection_data.collect()
    show_begin_points()
    begin_point = check_user_answer()
    removed = processing.process(values, begin_point)
    output = "COMPUTER: Has removed " + str(removed) + " members."
    print(output)
    
    main_menu()
