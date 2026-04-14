interface Props {
  children: React.ReactNode;
  onClick?: () => void;
  variant?: "primary" | "danger" | "secondary";
  disabled?: boolean;
}

export default function Button({
  children,
  onClick,
  variant = "primary",
  disabled,
}: Props) {
  return (
    <button
      onClick={onClick}
      disabled={disabled}
      className={`button ${variant}`}
    >
      {children}
    </button>
  );
}
