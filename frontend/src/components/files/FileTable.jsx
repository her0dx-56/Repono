import { useState } from 'react';
import {
  ArrowDownWideNarrow,
  File,
  FileArchive,
  FileAudio,
  FileImage,
  FileText,
  FileVideo,
} from 'lucide-react';
import { formatFileSize, formatDate } from '../../utils/formatFileSize.js';
import { fileCategory } from '../../utils/fileHelpers.js';
import FileActions from './FileActions.jsx';
import EmptyState from '../common/EmptyState.jsx';

const icons = {
  image: FileImage,
  video: FileVideo,
  audio: FileAudio,
  archive: FileArchive,
  pdf: FileText,
  text: FileText,
  file: File,
};

export default function FileTable({
  files,
  onDownload,
  onShare,
  onDelete,
  onPreview,
}) {
  const [sort, setSort] = useState('newest');

  const ordered = [...files].sort((a, b) => {
    if (sort === 'name') return a.name.localeCompare(b.name);
    if (sort === 'size') return (b.size || 0) - (a.size || 0);
    const timeB = new Date(b.modified || b.createdAt || 0).getTime();
    const timeA = new Date(a.modified || a.createdAt || 0).getTime();
    return timeB - timeA;
  });

  if (!files.length) {
    return (
      <EmptyState
        title="No files found"
        body="Upload a file or try a different search."
      />
    );
  }

  return (
    <div className="table-wrap">
      <table className="file-table">
        <thead>
          <tr>
            <th>FILE NAME</th>
            <th>
              <button
                className="sort-button"
                onClick={() => setSort(sort === 'size' ? 'newest' : 'size')}
              >
                SIZE <ArrowDownWideNarrow size={13} />
              </button>
            </th>
            <th>LAST MODIFIED</th>
            <th>ACCESS</th>
            <th />
          </tr>
        </thead>
        <tbody>
          {ordered.map((f) => {
            const I = icons[fileCategory(f.type, f.name)] || File;
            const ext =
              f.name && f.name.includes('.')
                ? f.name.split('.').pop().toUpperCase()
                : f.type?.split('/')[1]?.toUpperCase() || 'FILE';

            return (
              <tr key={f.id}>
                <td>
                  <button className="table-file" onClick={() => onPreview(f)}>
                    <span className={`file-type ${fileCategory(f.type, f.name)}`}>
                      <I size={20} />
                    </span>
                    <span>
                      <b>{f.name}</b>
                      <small>{ext}</small>
                    </span>
                  </button>
                </td>
                <td>{formatFileSize(f.size)}</td>
                <td>{formatDate(f.modified || f.createdAt)}</td>
                <td>
                  <span className={f.shared ? 'status-pill shared' : 'status-pill'}>
                    {f.shared ? 'Shared' : 'Private'}
                  </span>
                </td>
                <td>
                  <FileActions
                    file={f}
                    onDownload={onDownload}
                    onShare={onShare}
                    onDelete={onDelete}
                  />
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}
