import React, { useEffect } from "react";
import { Message, User } from "@/types/playerTypes";
import Image from "next/image";
import { getDayFromTimestamp } from "@/utils/timeHelper";
import { classNames } from "@/utils/utilhelper";

type ChatTimelineProps = {
  messageLogList: Message[];
  userData: User;
};

const ChatTimeline: React.FC<ChatTimelineProps> = ({
  messageLogList,
  userData,
}) => {
  useEffect(() => {
    const container = window.document.getElementById("messageLogList");
    if (!container) return;
    container.scrollTop = container.scrollHeight;
  }, [messageLogList]);

  return (
    <div className="flex-1 min-h-[70vh] max-h-[70vh]">
      <h3 className="mt-4 font-bold text-xl">Chat</h3>
      <ul className="overflow-y-scroll max-h-[60vh] my-4" id="messageLogList">
        {messageLogList &&
          messageLogList.length > 0 &&
          messageLogList.map((messageLog: Message, idx) => {
            if (messageLog.message_type === "system") return;
            return (
              <li key={idx} className="flex py-0.5 pr-16 leading-[22px]">
                <div className="flex py-1 leading-[22px]">
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
                        className={`mr-2 ${classNames(
                          messageLog.user.username === userData.username
                            ? "font-bold text-[#3b82f6]"
                            : messageLog.user.username === "server"
                            ? "font-bold text-orange-500"
                            : "font-normal text-black"
                        )}`}
                      >
                        {messageLog.user.username}
                      </span>
                      <span className="text-xs font-medium text-gray-900">
                        {getDayFromTimestamp(messageLog.timestamp)}
                      </span>
                    </p>
                    <p
                      className={`${classNames(
                        messageLog.user.username === "server"
                          ? "italic text-gray-500 font-bold"
                          : "text-gray-900"
                      )}`}
                    >
                      {messageLog.message}
                    </p>
                  </div>
                </div>
              </li>
            );
          })}
      </ul>
    </div>
  );
};

export default ChatTimeline;
