# coding: utf-8

import vk_api

def starter():
    vk_session = authorization()
    location_id = input_location_id()
    my_id = input_my_id()

    posts_offset = 0
    while True:
        values = {
            "owner_id": location_id,
            "count": 100,
            "offset": posts_offset,
            "filter": "all",
            "v": 5.126
        }

        posts = get_posts(vk_session, values)
        i = 0
        while i < len(posts):
            comments_offset = 0
            while True:
                values = {
                    "owner_id": posts[i]["owner_id"],
                    "post_id": posts[i]["id"],
                    "count": 100,
                    "offset": comments_offset,
                    "v": 5.68
                }
                comments = get_comments(vk_session, values)
                my_comments = select_my_comments(my_id, comments)

                n = 0
                while n < len(my_comments):
                    values = {
                        "owner_id": location_id,
                        "comment_id": my_comments[n]["id"],
                        "v": 5.126
                    }
                    delete_comment(vk_session, posts[i]["id"], values)
                    n += 1
                if len(comments) >= 100:
                    comments_offset += 100
                else:
                    break

            i += 1
        if len(posts) >= 100:
            posts_offset += 100
            print("Current offset = " + str(posts_offset))
        else:
            print("All is done!")
            break


def authorization():
    print("Enter access token:")
    access_token = raw_input("> ")

    vk_session = vk_api.VkApi(token=access_token)
    vk_session._auth_token()
    return vk_session


def input_location_id():
    print("Enter ID of location (community or user's page):")
    location_id = int(raw_input("> "))
    return location_id


def input_my_id():
    print("Enter your ID (or ID of author of comments to be removed):")
    my_id = int(raw_input("> "))
    return my_id


def get_posts(vk_session, values):
    response = vk_session.method("wall.get", values)
    return response["items"]


def get_comments(vk_session, values):
    response = vk_session.method("wall.getComments", values)
    return response["items"]


def select_my_comments(my_id, comments):
    my_comments = []
    i = 0
    while i < len(comments):
        if comments[i]["from_id"] == my_id:
            my_comments.append(comments[i])
        i += 1
    return my_comments


def delete_comment(vk_session, post_id, values):
    response = vk_session.method("wall.deleteComment", values)
    if response == 1:
        print("Comment https://vk.com/wall" + str(values["owner_id"]) + "_" + str(post_id) + "?reply=" + str(values["comment_id"]) + " has been deleted.")
    else:
        print("https://vk.com/wall" + str(values["owner_id"]) + "_" + str(post_id) + "?reply=" + str(values["comment_id"]))
        print(response)


starter()