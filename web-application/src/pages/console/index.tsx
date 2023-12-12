import AdminMessageSender from "@/components/settings/AdminMessageSender";
import SetMaxPayer from "@/components/settings/SetMaxPlayer";
import ClearMessageLogList from "@/components/settings/clearMessageLogList";
import useWebSocket from "@/hook/useWebSocket";
import { adminUser } from "@/store/chatStore";

export default function Console() {
  const {
    socket,
    isConnected,
    userList,
    user,
    messageLogList,
    handleWebSocketMessageSend,
  } = useWebSocket(null, adminUser);

  const fullUserList = [adminUser, ...userList];
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
    </main>
  );
}
