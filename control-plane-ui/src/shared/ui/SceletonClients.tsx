export const SkeletonClients = () => (
  <div style={{ display: "grid", gap: 12 }}>
    {[...Array(5)].map((_, i) => (
      <div key={i} style={{
        height: 80,
        background: "#eee",
        borderRadius: 12
      }} />
    ))}
  </div>
);
