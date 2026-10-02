export function fileCategory(type = '', name = '') {
  const combined = `${type} ${name}`.toLowerCase();
  if (combined.includes('image') || /\.(png|jpe?g|gif|webp|svg|bmp|ico)$/i.test(name)) return 'image';
  if (combined.includes('video') || /\.(mp4|webm|mkv|mov|avi)$/i.test(name)) return 'video';
  if (combined.includes('audio') || /\.(mp3|wav|ogg|m4a|flac)$/i.test(name)) return 'audio';
  if (combined.includes('pdf') || /\.pdf$/i.test(name)) return 'pdf';
  if (combined.includes('zip') || combined.includes('compressed') || /\.(zip|tar|gz|rar|7z)$/i.test(name)) return 'archive';
  if (combined.includes('text') || combined.includes('json') || /\.(txt|md|csv|json|js|jsx|ts|tsx|html|css|go|py)$/i.test(name)) return 'text';
  return 'file';
}

export function downloadBlob(blob, name) {
  const fileName = name || blob?.name || 'download';
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = fileName;
  document.body.appendChild(a);
  a.click();
  a.remove();
  setTimeout(() => URL.revokeObjectURL(url), 1000);
}
