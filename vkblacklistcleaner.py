# config: utf-8


import vk_api
import time


def exception_handler(sender, var_except):
    try:
        if str(var_except).lower().find("captcha needed") != -1:
            print(
                "COMPUTER [" + sender + "]: Error, " +
                str(var_except) + ". " +
                "Timeout: 60 sec.")
            time.sleep(60)

            return

        elif str(var_except).lower().find("failed to establish " +
                                          "a new connection") != -1:
            print(
                "COMPUTER [" + sender + "]: Error, " +
                str(var_except) + ". " +
                "Timeout: 60 sec.")
            time.sleep(60)

            return

        elif str(var_except).lower().find("connection aborted") != -1:
            print(
                "COMPUTER [" + sender + "]: Error, " +
                str(var_except) + ". " +
                "Timeout: 60 sec.")
            time.sleep(60)

            return

        else:
            print(
                "COMPUTER [" + sender + "]: Error, " +
                str(var_except) +
                ". Exit from program...")
            exit(0)
    except Exception as var_except:
        sender += " -> Exception handler"
        print(
            "COMPUTER [" + sender + "]: Error, " +
            str(var_except) +
            ". Exit from program...")
        exit(0)


def starter():
    sender = "Starter"

    try:
        user_answer = raw_input("USER [" + sender + " -> Get token]: ")
        access_token = user_answer
        vk_session = vk_api.VkApi(token=access_token)
        vk_session._auth_token()

        user_answer = raw_input("USER [" + sender + " -> Get group id]: ")

        group_id = int(user_answer)

        main(sender, vk_session, group_id)

    except Exception as var_except:
        exception_handler(sender, var_except)
        return starter()


def unban_users(sender, vk_session, group_id):
    sender += " -> Unban users"

    def get_user_list(sender, vk_session, get_user_list, offset):
        sender += " -> Get user list"

        try:
            values = {
                "group_id": group_id,
                "count": 200,
                "offset": offset
            }

            response = vk_session.method("groups.getBanned", values)

            banned_list = response["items"]

            return banned_list

        except Exception as var_except:
            exception_handler(sender, var_except)
            unban_users(sender, vk_session, group_id)

    def get_user_name(sender, vk_session, item):
        sender += " -> Get user's name"

        try:

            author_values = {
                "user_ids": item["profile"]["id"]
            }

            response = vk_session.method("users.get",
                                         author_values)

            first_name = response[0]["first_name"]
            last_name = response[0]["last_name"]

            user_name = first_name + " " + last_name

            return user_name

        except Exception as var_except:
            exception_handler(sender, var_except)
            return get_user_name(sender, vk_session, item)

    def get_group_name(sender, vk_session, item):
        sender += " -> Get user's name"

        try:

            author_values = {
                "group_ids": item["group"]["id"]
            }

            response = vk_session.method("groups.getById", author_values)

            group_name = response[0]["name"]

            return group_name

        except Exception as var_except:
            exception_handler(sender, var_except)
            return get_group_name(sender, vk_session, item)

    def run_unban(sender, vk_session, group_id, item):
        sender += " -> Run unban"

        try:
            if item["type"] == "profile":
                values = {
                    "group_id": group_id,
                    "owner_id": item["profile"]["id"]
                }
            else:
                values = {
                    "group_id": group_id,
                    "owner_id": "-" + str(item["group"]["id"])
                }

            vk_session.method("groups.unban", values)

        except Exception as var_except:
            if str(var_except).lower().find("user not blacklisted") != -1:
                return
            exception_handler(sender, var_except)
            return run_unban(sender, vk_session, group_id, item)

    try:

        banned_list = []
        offset = 0
        while True:
            else_banned_list = get_user_list(sender, vk_session, group_id, offset)
            banned_list.extend(else_banned_list)
            if len(else_banned_list) < 200:
                break
            offset += 200

        print("COMPUTER [" + sender + "]: Black list contained " +
              str(len(banned_list)) + " subjects.")

        i = len(banned_list) - 1

        removed = 0

        while i >= 0:

            item = banned_list[i]

            if "profile" in item:
                if "id" in item["profile"]:

                    user_name = get_user_name(sender, vk_session, item)
                    run_unban(sender, vk_session, group_id, item)

                    print(str(i + 1) + ". " + user_name +
                        " has been removed from black list.")

                    removed += 1
            if "group" in item:
                if "id" in item["group"]:

                    group_name = get_group_name(sender, vk_session, item)
                    run_unban(sender, vk_session, group_id, item)

                    print(str(i + 1) + ". " + group_name +
                        " has been removed from black list.")

                    removed += 1

            i -= 1

        return removed

    except Exception as var_except:
        exception_handler(sender, var_except)
        return unban_users(sender, vk_session, group_id)


def main(sender, vk_session, group_id):
    sender += " -> Main"

    try:

        removed = unban_users(sender, vk_session, group_id)

        print(str(removed) + " users has been removed from blacklist.")

    except Exception as var_except:
        exception_handler(sender, var_except)
        return main(sender, vk_session, group_id)


starter()
