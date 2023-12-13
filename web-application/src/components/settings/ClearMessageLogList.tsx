import { adminUser, initialUserSelection } from "@/store/chatStore";
import { maxPlayerAtom } from "@/store/gameAtom";
import { Message } from "@/types/playerTypes";
import { WsJsonRequest } from "@/types/wsTypes";
import { useAtom, useAtomValue } from "jotai";
import { useState } from "react";
type ClearMessageLogListProps = {
  messageLogList: Message[];
  sendMessage: (message: WsJsonRequest) => void;
};
export default function ClearMessageLogList({
  messageLogList,
  sendMessage,
}: ClearMessageLogListProps) {
  const maxPlayer = useAtomValue(maxPlayerAtom);

  const handleClearMessageLogList = (
    event: React.MouseEvent<HTMLButtonElement>
  ) => {
    event.preventDefault();
    const jsonData: WsJsonRequest = {
      max_player: maxPlayer,
      action: "clear_messages",
      timestamp: Date.now(),
      user: adminUser,
      message: `Clear messages`,
      game_round: 0,
      game_turns_left: 0,
      user_selection: initialUserSelection,
    };
    sendMessage(jsonData);
  };

  return (
    <>
      <h3 className="mt-4 mb-2 font-bold text-xl">
        Clear message log list - current num of message is{" "}
        <span className="text-[#3b82f6]">
          {messageLogList && messageLogList.length}
        </span>
      </h3>
      <button
        onClick={handleClearMessageLogList}
        className="px-2 py-1 text-sm font-medium text-white bg-gray-900 rounded-md hover:bg-[#3b82f6] focus:outline-none focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75"
      >
        Clear Message Log List
      </button>
    </>
  );
}
