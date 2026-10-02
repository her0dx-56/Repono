import {
  File,
  FileArchive,
  FileAudio,
  FileImage,
  FileText,
  FileVideo,
  Download,
  Share2,
  Trash2,
} from 'lucide-react';
import { formatFileSize } from '../../utils/formatFileSize.js';
import { fileCategory } from '../../utils/fileHelpers.js';

const icons = {
  image: FileImage,
  video: FileVideo,
  audio: FileAudio,
  archive: FileArchive,
  pdf: FileText,
  text: FileText,
  file: File,
};

export default function FileCard({
  file,
  onDownload,
  onShare,
  onDelete,
  onPreview,
}) {
  const I = icons[fileCategory(file.type, file.name)] || File;
  const dateStr =
    file.modified || file.createdAt
      ? new Date(file.modified || file.createdAt).toLocaleDateString()
      : 'Today';

  return (
    <article className="file-card">
      <button
        className={`file-type ${fileCategory(file.type, file.name)}`}
        onClick={() => onPreview(file)}
        aria-label={`Details for ${file.name}`}
      >
        <I size={23} />
      </button>

      <button className="file-card-name" onClick={() => onPreview(file)}>
        <b>{file.name}</b>
        <small>
          {formatFileSize(file.size)} · {dateStr}
        </small>
      </button>

      <span className={file.shared ? 'status-pill shared' : 'status-pill'}>
        {file.shared ? 'Shared' : 'Private'}
      </span>

      <div className="card-actions">
        <button
          className="icon-button"
          title="Download"
          onClick={() => onDownload(file)}
        >
          <Download size={17} />
        </button>

        {!file.isSharedWithMe && (
          <button
            className="icon-button"
            title="Share"
            onClick={() => onShare(file)}
          >
            <Share2 size={17} />
          </button>
        )}

        {!file.isSharedWithMe && (
          <button
            className="icon-button danger-hover"
            title="Delete"
            onClick={() => onDelete(file)}
          >
            <Trash2 size={17} />
          </button>
        )}
      </div>
    </article>
  );
}
