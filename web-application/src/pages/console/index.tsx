import ChatTimeline from "@/components/ChatTimeline";
import AdminMessageSender from "@/components/settings/AdminMessageSender";
import ClearMessageLogList from "@/components/settings/ClearMessageLogList";
import RestartGame from "@/components/settings/RestartGame";
import RestartRound from "@/components/settings/RestartRound";
import SetGameNums from "@/components/settings/SetGameNums";
import SetMaxPayer from "@/components/settings/SetMaxPlayer";
import useWebSocket from "@/hook/useWebSocket";
import { messageLogListAtom } from "@/store/chatAtom";
import { adminUser } from "@/store/chatStore";
import {
  gameRoundAtom,
  userSelectionAtom,
  userSelectionListAtom,
} from "@/store/gameAtom";
import { UserSelection } from "@/types/wsTypes";
import { useAtomValue } from "jotai";

export default function Console() {
  const {
    socket,
    isConnected,
    userList,
    user,
    messageLogList,
    handleWebSocketMessageSend,
  } = useWebSocket(null, adminUser);

  const userSelectionList = useAtomValue(userSelectionListAtom);
  const fullUserList = [adminUser, ...userList];
  const gameRound = useAtomValue(gameRoundAtom);
  // const fullUserList = userList.concat(user);
  return (
    <main className="py-8 mx-auto w-[80vw] max-w-2xl min-h-max relative">
      <h1 className="text-3xl font-bold mb-4">Console</h1>
      <SetMaxPayer sendMessage={handleWebSocketMessageSend} />
      <AdminMessageSender
        userList={fullUserList}
        sendMessage={handleWebSocketMessageSend}
      />
      <ClearMessageLogList
        messageLogList={messageLogList}
        sendMessage={handleWebSocketMessageSend}
      />
      <SetGameNums sendMessage={handleWebSocketMessageSend} />

      <br className="my-4" />
      <RestartRound sendMessage={handleWebSocketMessageSend} />
      <RestartGame sendMessage={handleWebSocketMessageSend} />
      <br className="my-4" />

      <h3 className="mt-4 mb-2 font-bold text-xl">Roud {gameRound}</h3>
      <ChatTimeline messageLogList={messageLogList} userData={adminUser} />
      {userSelectionList &&
        userSelectionList.length > 0 &&
        userSelectionList.map((selection: UserSelection, index) => (
          <div key={index} className="flex flex-col mb-4">
            <p className="mb-2 font-bold">
              {selection.user.username} Selection:{" "}
              <span className="text-[#3b82f6]">{selection.selection}</span>
            </p>
            <p className="">
              {selection.user.username} Reason:{" "}
              <span className="text-[#3b82f6]">{selection.reason}</span>
            </p>
          </div>
        ))}
    </main>
  );
}
