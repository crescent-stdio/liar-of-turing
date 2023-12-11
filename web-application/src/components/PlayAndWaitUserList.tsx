import { User } from "@/types/playerTypes";
import OnlineUserList from "./OnlineUserList";
type PlayAndWaitUserListProps = {
  userList: User[];
  userData: User;
};
export default function PlayAndWaitUserList({
  userList,
  userData,
}: PlayAndWaitUserListProps) {
  return (
    <div>
      <OnlineUserList userData={userData} userList={userList} />
    </div>
  );
}
