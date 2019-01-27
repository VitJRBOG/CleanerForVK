# coding: utf-8
u"""Модуль ввода данных пользователем."""


def get_access_token():
    u"""Получает токен доступа."""
    user_answer = raw_input("USER [.. -> New token]: ")
    return user_answer


def get_public_url():
    u"""Получает URL сообщества."""
    user_answer = raw_input("USER [.. -> URL]: ")
    return user_answer


def get_number_members():
    u"""Получает количество пользователей."""
    user_answer = raw_input("USER [.. -> Number members]: ")
    return user_answer


def get_number_members_to_remove():
    u"""Получает количество пользователей для удаления."""
    user_answer = raw_input("USER [.. -> Members to remove]: ")
    return user_answer


def get_sort():
    u"""Получает метод сортировки пользователей в источнике."""
    user_answer = raw_input("USER [.. -> Begin from]: ")
    return user_answer
