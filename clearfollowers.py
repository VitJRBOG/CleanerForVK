# coding: utf-8


import vk_api


def main():
    sender = "Main"
    user_answer = raw_input("USER [" + sender + " -> New token]: ")

    vk_session = authorize(sender, user_answer)

    items = get_followers(sender, vk_session)

    i = 0

    while i < len(items):

        user_id = items[i]

        item = get_user(sender, vk_session, user_id)

        if "deactivated" in item[0]:
            if item[0]["deactivated"] == "deleted" or\
               item[0]["deactivated"] == "banned":
                ban_user(sender, vk_session, item[0])

        i += 1


def authorize(sender, access_token):
    sender += " -> Authorize"

    vk_session = vk_api.VkApi(token=access_token)
    vk_session._auth_token()

    return vk_session


def get_followers(sender, vk_session):
    sender += " -> Get followers"

    values = {
        "count": 1000
    }

    response = vk_session.method("users.getFollowers", values)

    return response["items"]


def get_user(sender, vk_session, item):
    sender += " -> Get user"

    values = {
        "user_ids": item
    }

    response = vk_session.method("users.get", values)

    return response


def ban_user(sender, vk_session, user):

    sender += " -> Ban user"

    values = {
        "owner_id": user["id"]
    }

    vk_session.method("account.ban", values)

    print(str(user["id"]) + " added to blacklist.")

main()
