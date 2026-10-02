import { useRef, useState } from 'react';
import { CloudUpload, FileUp } from 'lucide-react';

export default function UploadZone({ onFiles }) {
  const input = useRef(null);
  const [drag, setDrag] = useState(false);

  function accept(list) {
    const files = [...list];
    if (files.length) onFiles(files);
  }

  return (
    <div
      className={`upload-zone ${drag ? 'dragging' : ''}`}
      onDragOver={(e) => {
        e.preventDefault();
        setDrag(true);
      }}
      onDragLeave={() => setDrag(false)}
      onDrop={(e) => {
        e.preventDefault();
        setDrag(false);
        accept(e.dataTransfer.files);
      }}
    >
      <input
        ref={input}
        type="file"
        multiple
        hidden
        onChange={(e) => {
          accept(e.target.files);
          e.target.value = '';
        }}
      />
      <div className="upload-icon">
        <CloudUpload size={28} />
      </div>
      <div>
        <b>Drop files here to send them on their way</b>
        <p>
          or{' '}
          <button
            className="text-button"
            onClick={() => input.current?.click()}
          >
            browse from your device
          </button>
        </p>
        <small>Any file type · Direct upload to Repono storage</small>
      </div>
      <button
        className="btn btn-outline browse-btn"
        onClick={() => input.current?.click()}
      >
        <FileUp size={17} /> Choose files
      </button>
    </div>
  );
}
