export function formatFileSize(bytes = 0) {
  const n = Number(bytes) || 0;
  if (n < 1000) return `${n} B`;
  if (n < 1e6) return `${(n / 1000).toFixed(1)} KB`;
  if (n < 1e9) return `${(n / 1e6).toFixed(2)} MB`;
  return `${(n / 1e9).toFixed(2)} GB`;
}

export function formatDate(value) {
  if (!value) return "Just now";
  const date = new Date(value);
  if (isNaN(date.getTime())) return "Recently";
  return new Intl.DateTimeFormat(undefined, {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(date);
}
