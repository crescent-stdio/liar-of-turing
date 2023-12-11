import React from "react";
import { User } from "@/types/playerTypes";

type OnlineUserListProps = {
  userList: User[];
  userData: User;
};

const OnlineUserList: React.FC<OnlineUserListProps> = ({
  userList,
  userData,
}) => {
  return (
    <div className="flex flex-col ml-4">
      <h3 className="mt-6 font-bold text-xl">Now online</h3>
      <ul className="my-4 h-80 flex-1">
        {userList &&
          userList.length > 0 &&
          userList.map((user: User, index) => {
            return (
              <li
                key={index}
                style={{
                  color:
                    user.username === userData.username ? "#3b82f6" : "black",
                  fontWeight: user.username === userData.username ? 700 : 400,
                }}
              >
                {user.username}
              </li>
            );
          })}
      </ul>
    </div>
  );
};

export default OnlineUserList;
