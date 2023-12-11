import React from "react";
import { Message, User } from "@/types/playerTypes";
import Image from "next/image";
import { getDayFromTimestamp } from "@/utils/timeHelper";

type ChatTimelineProps = {
  messageLogList: Message[];
  userData: User;
};

const ChatTimeline: React.FC<ChatTimelineProps> = ({
  messageLogList,
  userData,
}) => {
  return (
    <ul
      className="overflow-y-scroll min-h-[60vh] max-h-[60vh] my-4 flex-1"
      id="messageLogList"
    >
      {messageLogList &&
        messageLogList.length > 0 &&
        messageLogList.map((messageLog: Message, idx) => {
          if (messageLog.message_type === "system") return;
          return (
            <li key={idx} className="flex py-0.5 pr-16 leading-[22px]">
              <div className="flex py-0.5 leading-[22px]">
                <div className="overflow-hidden relative mt-0.5 mr-2 w-10 min-w-fit h-10 rounded-sm">
                  <Image
                    src={`/nickname_icon/${messageLog.user.nickname_id}.png`}
                    alt={`${messageLog.user.username} icon`}
                    width={40}
                    height={40}
                  />
                </div>
                <div>
                  <p className="flex items-baseline">
                    <span
                      className="mr-2 font-bold text-green-400"
                      style={{
                        color:
                          messageLog.user.username === userData.username
                            ? "#3b82f6"
                            : "black",
                      }}
                    >
                      {messageLog.user.username}
                    </span>
                    <span className="text-xs font-medium text-gray-900">
                      {getDayFromTimestamp(messageLog.timestamp)}
                    </span>
                  </p>
                  <p className="text-gray-900">{messageLog.message}</p>
                </div>
              </div>
            </li>
          );
        })}
    </ul>
  );
};

export default ChatTimeline;
