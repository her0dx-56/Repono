
import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import {
  ArrowLeft,
  ArrowDownToLine,
  Cloud,
  FileText,
  HardDrive,
  ShieldCheck,
  AlertCircle,
} from "lucide-react";

import Brand from "../components/common/Brand.jsx";
import { fileService } from "../services/fileService.js";
import { formatFileSize } from "../utils/formatFileSize.js";

export default function SharedFilePage() {
  const { id } = useParams();
  const [fileInfo, setFileInfo] = useState(null);
  const [loading, setLoading] = useState(true);
  const [downloading, setDownloading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    let isMounted = true;

    async function loadFileInfo() {
      if (!id) {
        setError("Invalid share link.");
        setLoading(false);
        return;
      }

      setLoading(true);
      setError("");

      try {
        const info = await fileService.getPublicShareInfo(id);
        if (isMounted) {
          setFileInfo(info);
        }
      } catch (err) {
        if (isMounted) {
          console.warn("Could not load public share info:", err);
          setError(
            err.message || "This shared file is unavailable or the link has expired."
          );
        }
      } finally {
        if (isMounted) {
          setLoading(false);
        }
      }
    }

    loadFileInfo();

    return () => {
      isMounted = false;
    };
  }, [id]);

  function download() {
    if (!id || downloading) return;

    try {
      setDownloading(true);

      const downloadUrl = fileService.getPublicDownloadUrl(id);
      const link = document.createElement("a");
      link.href = downloadUrl;
      link.download = "";
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);

      setTimeout(() => {
        setDownloading(false);
      }, 2000);
    } catch (err) {
      console.error("Failed to start download:", err);
      setDownloading(false);
    }
  }

  const fileName = fileInfo?.name || "Shared file";
  const fileSize = fileInfo?.size || 0;

  return (
    <div className="public-page">
      <header>
        <Link to="/home">
          <Brand />
        </Link>

        <Link className="btn btn-outline" to="/home">
          <ArrowLeft size={16} /> Back to Repono
        </Link>
      </header>

      <main className="public-card">
        <div className="public-art">
          <div className="public-sun" />
          <div className="public-cloud">
            <Cloud size={78} />
          </div>
          <div className="public-file">
            <FileText size={32} />
          </div>
        </div>

        {error ? (
          <>
            <span className="unavailable-icon">
              <AlertCircle />
            </span>

            <p className="eyebrow">SHARE LINK UNAVAILABLE</p>
            <h1>Unable to access this file.</h1>
            <p className="public-description">{error}</p>

            <Link to="/home" className="btn btn-red btn-large">
              Back to Repono
            </Link>
          </>
        ) : loading ? (
          <>
            <p className="eyebrow">A FILE WAS SHARED WITH YOU</p>
            <h1>Loading shared file...</h1>
            <p className="public-description">
              Retrieving file details from Repono cloud...
            </p>
          </>
        ) : (
          <>
            <p className="eyebrow">A FILE WAS SHARED WITH YOU</p>
            <h1>{fileName}</h1>

            <p className="public-description">
              This file was shared with you through Repono.
              Use the button below to download it securely to your device.
            </p>

            <div
              className="public-meta"
              style={{
                flexWrap: "wrap",
                gap: "10px 16px",
                maxWidth: "100%",
              }}
            >
              <span
                title={fileName}
                style={{
                  maxWidth: "220px",
                  minWidth: 0,
                }}
              >
                <FileText size={16} style={{ flexShrink: 0 }} />
                <span
                  style={{
                    overflow: "hidden",
                    textOverflow: "ellipsis",
                    whiteSpace: "nowrap",
                  }}
                >
                  {fileName}
                </span>
              </span>

              {fileSize > 0 && (
                <span style={{ flexShrink: 0 }}>
                  <HardDrive size={16} style={{ flexShrink: 0 }} />
                  {formatFileSize(fileSize)}
                </span>
              )}

              <span style={{ flexShrink: 0 }}>
                <ShieldCheck size={16} style={{ flexShrink: 0 }} />
                Public share link
              </span>
            </div>

            <button
              className="btn btn-red btn-large"
              onClick={download}
              disabled={downloading}
            >
              <ArrowDownToLine size={18} />
              {downloading ? "Downloading..." : "Download file"}
            </button>
          </>
        )}
      </main>

      <footer>REPONO · YOUR FILES. YOUR WORLD.</footer>
    </div>
  );
}