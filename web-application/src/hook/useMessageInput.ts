import { messageAtom } from "@/store/chatAtom";
import { useAtom } from "jotai";
import { useState } from "react";

export interface MessageInputHook {
  message: string;
  handleMessageChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
  handleSubmit: (event: React.FormEvent<HTMLFormElement>) => void;
  resetMessage: () => void;
}

export default function useMessageInput(): MessageInputHook {
  const [message, setMessage] = useAtom(messageAtom);

  // Handles the change in the input field
  const handleMessageChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setMessage(event.target.value);
  };

  // Handles the message submission
  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (message.trim() !== "") {
      setMessage(message.trim());
      // resetMessage();
    }
  };

  // Resets the input field
  const resetMessage = () => {
    setMessage("");
  };

  return {
    message,
    handleMessageChange,
    handleSubmit,
    resetMessage,
  };
}
