# coding: utf8
u"""Модуль получения предварительных данных."""


import vk_api
import input_data


def collect():
    u"""Запускает функции сборки данных."""
    access_token = input_data.get_access_token()
    vk_session = authorize(access_token)
    public_url = input_data.get_public_url()
    public_id = select_public_id(vk_session, public_url)
    number_members = int(input_data.get_number_members())
    number_members_to_remove = int(input_data.get_number_members_to_remove())
    values = {
        "vk_session": vk_session,
        "public_id": public_id,
        "number_members": number_members,
        "number_to_remove": number_members_to_remove
    }
    return values


def authorize(access_token):
    vk_session = vk_api.VkApi(token=access_token)
    vk_session._auth_token()

    return vk_session


def select_public_id(vk_session, public_url):
    u"""Получает id сообщества ВК."""
    def select_public_domain(public_url):
        u"""Извлекает домайн из ссылки на сообщество ВК."""
        underline_for_find = "vk.com/"
        bgn_indx = public_url.find(underline_for_find)
        end_indx = bgn_indx + len(underline_for_find)
        public_domain = public_url[end_indx:]

        return public_domain

    def get_id(vk_session, public_domain):
        u"""Запрашивает id сообщества ВК."""
        values = {
            "group_id": public_domain
        }
        response = vk_session.method("groups.getById", values)
        public_id = str(response[0]["id"])

        return public_id
    
    public_domain = select_public_domain(public_url)
    public_id = get_id(vk_session, public_domain)
    return public_id
