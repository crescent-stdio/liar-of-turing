import { User } from "@/types/playerTypes";
import OnlineUserList from "./OnlineUserList";
import { playerListAtom, watcherListAtom } from "@/store/chatAtom";
import { useAtomValue } from "jotai";
import { maxPlayerAtom } from "@/store/gameAtom";
import { useEffect, useState } from "react";
type PlayAndWaitUserListProps = {
  userData: User;
};
export default function PlayAndWaitUserList({
  userData,
}: PlayAndWaitUserListProps) {
  const playerList = useAtomValue(playerListAtom);
  const watcherList = useAtomValue(watcherListAtom);
  const maxPlayer = useAtomValue(maxPlayerAtom);
  const [playerListTitle, setPlayerListTitle] = useState<string>("Players");

  // const playerListTitle = `Players [${playerList.length}/${maxPlayer}]`;
  console.log(maxPlayer);
  useEffect(() => {
    setPlayerListTitle(`Players [${playerList.length}/${maxPlayer}]`);
  }, [playerList.length, maxPlayer]);
  return (
    <div className="flex flex-col h-full ml-4">
      <OnlineUserList
        name={playerListTitle}
        userData={userData}
        userList={playerList}
      />
      <hr className="my-2" />
      <OnlineUserList
        name={"Watchers"}
        userData={userData}
        userList={watcherList}
      />
    </div>
  );
}
