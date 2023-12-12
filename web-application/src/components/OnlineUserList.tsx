import React from "react";
import { User } from "@/types/playerTypes";
import { classNames } from "@/utils/utilhelper";

type OnlineUserListProps = {
  name: string;
  userList: User[];
  userData: User;
};

const OnlineUserList: React.FC<OnlineUserListProps> = ({
  name,
  userList,
  userData,
}) => {
  return (
    <div className="flex-1">
      <h3 className="font-bold text-xl">{name}</h3>
      <ul className="my-4 w-full md:min-h-80 flex-1 min-h-full">
        {userList &&
          userList.length > 0 &&
          userList.map((user: User, index) => {
            return (
              <li
                key={index}
                className={`${classNames(
                  "",
                  user.username === userData.username
                    ? "font-bold text-[#3b82f6]"
                    : "font-normal text-black"
                    ? name === "Watchers"
                      ? "text-gray-500"
                      : "text-black"
                    : "text-gray-500"
                )}`}
              >
                {user.username === userData.username
                  ? `${user.username} ★`
                  : user.username}
              </li>
            );
          })}
      </ul>
    </div>
  );
};

export default OnlineUserList;
