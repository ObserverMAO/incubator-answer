import { FC, MouseEvent, ReactNode, useEffect, useState } from 'react';
import { Modal } from 'react-bootstrap';

import classnames from 'classnames';

const VideoViewer: FC<{
  children: ReactNode;
  className?: classnames.Argument;
}> = ({ children, className }) => {
  const [visible, setVisible] = useState(false);
  const [videoSrc, setVideoSrc] = useState('');

  const onClose = () => {
    setVisible(false);
    setVideoSrc('');
  };

  const isVideoURL = (url: string) => {
    const videoExtensions = ['.mp4', '.mov', '.avi', '.webm', '.ogg'];
    return videoExtensions.some((ext) => url.toLowerCase().endsWith(ext));
  };

  const checkClickForVideoView = (evt: MouseEvent<HTMLElement>) => {
    const { target } = evt;
    if (target instanceof HTMLElement) {
      const videoElement = target.closest('video');
      if (videoElement) {
        const src = videoElement.currentSrc || videoElement.src;
        if (src) {
          setVideoSrc(src);
          setVisible(true);
        }
      } else if (target.tagName === 'A') {
        const href = target.getAttribute('href');
        if (href && isVideoURL(href)) {
          setVideoSrc(href);
          setVisible(true);
          evt.preventDefault();
        }
      }
    }
  };

  useEffect(() => {
    return () => {
      onClose();
    };
  }, []);

  return (
    <div
      className={classnames('video-viewer', className)}
      onClick={checkClickForVideoView}>
      {children}
      <Modal
        show={visible}
        fullscreen
        centered
        scrollable
        contentClassName="bg-transparent"
        onHide={onClose}>
        <Modal.Body onClick={onClose} className="video-viewer p-0 d-flex">
          <video
            className="w-100 h-100"
            controls
            autoPlay
            onClick={(evt) => {
              evt.stopPropagation();
            }}
            src={videoSrc}>
            <track kind="captions" />
          </video>
        </Modal.Body>
      </Modal>
    </div>
  );
};

export default VideoViewer;
