# coding: utf8
u"""Модуль обработки предварительных данных и отправки запросов."""


import output_data


def process(values, begin_point):
    u"""Запуск функций обработки."""
    number_to_remove = values["number_to_remove"]
    members = get_members(values, begin_point)
    vk_session = values["vk_session"]
    public_id = values["public_id"]
    condition = "banned"
    removed = 0
    for i, member in enumerate(members):
        user = get_user(vk_session, member)
        is_member = check_following(vk_session, public_id, user)
        if is_member:
            if "deactivated" in user:
                if user["deactivated"] == condition:
                    remove_member(vk_session, public_id, user)
                    removed += 1
                    if removed >= number_to_remove:
                        break
    return removed

def get_members(values, begin_point):
    u"""Получение списка подписчиков."""
    def algorithm_get_members(response_values, begin_point):
        u"""Алгоритм отправки запроса на получение списка подписчиков."""
        vk_session = response_values["vk_session"]
        public_id = response_values["public_id"]
        number_members = response_values["number_members"]
        offset = response_values["offset"]
        values = {
            "group_id": public_id,
            "count": number_members,
            "sort": begin_point,
            "offset": offset
        }
        response = vk_session.method("groups.getMembers", values)
        return response["items"]

    response_values = {
        "vk_session": values["vk_session"],
        "public_id": values["public_id"],
        "sort": begin_point,
        "offset": 0
    }

    members = []

    number_members = values["number_members"]
    if number_members > 1000:
        response_values.update({"number_members": 1000})
        left_number_members = number_members
        while left_number_members != 0:
            if left_number_members > 1000:
                items = algorithm_get_members(response_values, begin_point)
                if len(items) > 0:
                    members.extend(items)
                    response_values["offset"] += 1000
                    left_number_members -= 1000
                else:
                    break
            else:
                response_values["number_members"] = left_number_members
                items = algorithm_get_members(response_values, begin_point)
                if len(items) > 0:
                    members.extend(items)
                    left_number_members = 0
                else:
                    break
    else:
        response_values.update({"number_members": number_members})
        members.extend(algorithm_get_members(response_values, begin_point))

    output = "Has found " + str(len(members)) + " users."
    output_data.data_to_console(output)

    return members


def get_user(vk_session, item):
    u"""Получение данных о страничке пользователя."""
    values = {
        "user_ids": item
    }
    response = vk_session.method("users.get", values)
    return response[0]


def check_following(vk_session, public_id, user):
    u"""Проверка наличия подписки."""
    values = {
        "group_id": public_id,
        "user_id": user["id"]
    }
    response = vk_session.method("groups.isMember", values)
    if response == 1:
        is_member = True
    else:
        is_member = False

    return is_member


def remove_member(vk_session, public_id, user):
    u"""Удаление подписчика."""
    values = {
        "group_id": public_id,
        "user_id": user["id"]
    }
    vk_session.method("groups.removeUser", values)
    full_name = user["first_name"] + " " + user["last_name"]
    output = "(vk.com/id" + str(user["id"]) + ") " + full_name + \
        " removed from public."
    output_data.data_to_console(output)
