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
            "filter": "other",
            "v": 5.126
        }

        posts = get_posts(vk_session, values)
        my_posts = select_my_posts(my_id, posts)
        i = 0
        while i < len(my_posts):
            values = {
                "owner_id": my_posts[i]["owner_id"],
                "post_id": my_posts[i]["id"],
                "v": 5.126
            }
            delete_post(vk_session, values)
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
    print("Enter your ID (or ID of author of wallposts to be removed):")
    my_id = int(raw_input("> "))
    return my_id


def get_posts(vk_session, values):
    response = vk_session.method("wall.get", values)
    return response["items"]


def select_my_posts(my_id, posts):
    my_posts = []
    i = 0
    while i < len(posts):
        if posts[i]["from_id"] == my_id:
            my_posts.append(posts[i])
        i += 1
    return my_posts


def delete_post(vk_session, values):
    response = vk_session.method("wall.delete", values)
    if response == 1:
        print("Post https://vk.com/wall" + str(values["owner_id"]) + "_" + str(values["post_id"]) + " has been deleted.")
    else:
        print("https://vk.com/wall" + str(values["owner_id"]) + "_" + str(values["post_id"]))
        print(response)


starter()