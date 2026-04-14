import type React from "react";

interface Props {
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  placeholder?: string;
  type?: string;
  name?: string;
  required?: boolean;
  style?: React.CSSProperties;
  onKeyDown?: (e: React.KeyboardEvent<HTMLInputElement>) => void;
}

export default function Input(props: Props) {
  return (
    <input
      className="input"
      value={props.value}
      onChange={props.onChange}
      onKeyDown={props.onKeyDown}
      placeholder={props.placeholder}
      type={props.type || "text"}
    />
  );
}