import { Download, Share2, Trash2 } from 'lucide-react';
export default function FileActions({ file, onDownload, onShare, onDelete }) {
  return (
    <div className="file-actions">
      <button
        className="icon-button"
        title={`Download ${file.name}`}
        onClick={() => onDownload(file)}
      >
        <Download size={17} />
      </button>
      {!file.isSharedWithMe && (
        <button
          className="icon-button"
          title={`Share ${file.name}`}
          onClick={() => onShare(file)}
        >
          <Share2 size={17} />
        </button>
      )}
      {!file.isSharedWithMe && (
        <button
          className="icon-button danger-hover"
          title={`Delete ${file.name}`}
          onClick={() => onDelete(file)}
        >
          <Trash2 size={17} />
        </button>
      )}
    </div>
  );
}
