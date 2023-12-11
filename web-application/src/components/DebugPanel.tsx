import { useState } from "react";
import { User } from "@/types/playerTypes";
import { WsJsonRequest } from "@/types/wsTypes";

type DebugPanelProps = {
  userList: User[];
  sendMessage: (message: WsJsonRequest) => void;
};

const DebugPanel: React.FC<DebugPanelProps> = ({ userList, sendMessage }) => {
  const [isDebugMode, setIsDebugMode] = useState<boolean>(false);
  const [testUsername, setTestUsername] = useState<string>("");
  const [testMessage, setTestMessage] = useState<string>("");

  const handleTestSendMessage = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (testUsername && testMessage) {
      const testUserData = userList.find(
        (user: User) => user.username === testUsername
      );
      if (!testUserData) return;
      const jsonData: WsJsonRequest = {
        action: "new_message",
        user: testUserData,
        timestamp: Date.now(),
        message: testMessage,
      };
      sendMessage(jsonData);
      setTestUsername("");
      setTestMessage("");
    }
  };

  return (
    <>
      <button
        className="top-0 -right-[25vw] absolute text-white hover:text-black"
        onClick={() => setIsDebugMode((isDebugMode) => !isDebugMode)}
      >
        {isDebugMode ? "Hide" : "Show"} Debug Panel
      </button>
      {isDebugMode && (
        <div className="flex flex-col">
          <div className="mt-40"></div>
          <form className="flex flex-row" onSubmit={handleTestSendMessage}>
            <label htmlFor="username">Username</label>
            <select
              name="username"
              id="username"
              onChange={(e) => {
                setTestUsername(e.target.value);
              }}
              className="border-2 border-gray-400 rounded-md w-fit-content"
            >
              {userList &&
                userList.length > 0 &&
                userList.map((user: User, index) => {
                  return (
                    <option key={index} value={user.username}>
                      {user.username}
                    </option>
                  );
                })}
            </select>
            <label htmlFor="message" className="mx-2 ">
              Message
            </label>
            <input
              className="border-2 border-gray-400 rounded-md w-fit-content"
              type="text"
              id="message"
              value={testMessage}
              onChange={(e) => setTestMessage(e.target.value)}
            />
          </form>
        </div>
      )}
    </>
  );
};

export default DebugPanel;
