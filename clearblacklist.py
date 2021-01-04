# config: utf-8


import vk_api


def main():
    sender = "Main"

    user_answer = raw_input("USER [" + sender + " -> New token]: ")

    vk_session = authorize(sender, user_answer)

    banned_list = get_banned(sender, vk_session)

    i = 0
    while i < len(banned_list):

        item = banned_list[i]

        unban_user(sender, vk_session, item)

        i += 1


def authorize(sender, access_token):
    sender += " -> Authorize"

    vk_session = vk_api.VkApi(token=access_token)
    vk_session._auth_token()

    return vk_session


def get_banned(sender, vk_session):
    sender += " -> Get banned"

    values = {
        "count": 200
    }

    response = vk_session.method("account.getBanned", values)

    return response["items"]


def unban_user(sender, vk_session, item):
    sender += " -> Unban user"

    values = {
        "user_id": item["id"]
    }

    vk_session.method("account.unbanUser", values)

    print(str(item["id"]) + " removed from blacklist.")


main()
