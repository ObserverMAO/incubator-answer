/* eslint-disable @typescript-eslint/no-unused-vars */
import { useEffect, useState, memo } from 'react';
import { Button, Form, Modal, Tab, Tabs } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';

import { Modal as AnswerModal } from '@/components';
import ToolItem from '../toolItem';
import { IEditorContext, Editor } from '../types';
import { uploadImage } from '@/services';

let context: IEditorContext;
const Video = ({ editorInstance }) => {
  const [editor, setEditor] = useState<Editor>(editorInstance);
  const { t } = useTranslation('translation', { keyPrefix: 'editor' });

  const loadingText = `![${t('video.uploading')}...]()`;
  const item = {
    label: 'camera-video',
    keyMap: ['Ctrl-j'],
    tip: `${t('video.text')} (Ctrl+J)`,
  };
  const [currentTab, setCurrentTab] = useState('localVideo');
  const [visible, setVisible] = useState(false);
  const [link, setLink] = useState({
    value: '',
    isInvalid: false,
    errorMsg: '',
    type: '',
  });

  const [videoName, setVideoName] = useState({
    value: '',
    isInvalid: false,
    errorMsg: '',
  });

  const verifyVideoSize = (files: FileList) => {
    if (files.length === 0) {
      return false;
    }
    const filteredFiles = Array.from(files).filter(
      (file) => !file.type.startsWith('video/'),
    );

    if (filteredFiles.length > 0) {
      AnswerModal.confirm({
        content: t('video.form_video.fields.file.msg.only_video'),
      });
      return false;
    }
    const filteredVideos = Array.from(files).filter(
      (file) => file.size / 1024 / 1024 > 100,
    );

    if (filteredVideos.length > 0) {
      AnswerModal.confirm({
        content: t('video.form_video.fields.file.msg.max_size'),
      });
      return false;
    }
    return true;
  };

  const upload = (
    files: FileList,
  ): Promise<{ url: string; name: string }[]> => {
    const promises = Array.from(files).map(async (file) => {
      const url = await uploadImage({ file, type: 'post' });

      return {
        name: file.name,
        url,
      };
    });

    return Promise.all(promises);
  };
  const handleClick = () => {
    if (!link.value) {
      setLink({ ...link, isInvalid: true });
      return;
    }
    setLink({ ...link, type: '' });

    const text = `![${videoName.value}](${link.value})`;

    editor.replaceSelection(text);

    setVisible(false);

    editor.focus();
    setLink({ ...link, value: '' });
    setVideoName({ ...videoName, value: '' });
  };
  function dragenter(e) {
    e.stopPropagation();
    e.preventDefault();
  }

  function dragover(e) {
    e.stopPropagation();
    e.preventDefault();
  }

  const drop = async (e) => {
    const fileList = e.dataTransfer.files;

    const bool = verifyVideoSize(fileList);

    if (!bool) {
      return;
    }

    const startPos = editor.getCursor();

    const endPos = { ...startPos, ch: startPos.ch + loadingText.length };

    editor.replaceSelection(loadingText);
    editor.setReadOnly(true);
    const urls = await upload(fileList).catch((ex) => {
      console.error('upload file error: ', ex);
    });

    const text: string[] = [];
    if (Array.isArray(urls)) {
      urls.forEach(({ name, url }) => {
        if (name && url) {
          text.push(`[${name}](${url})`);
        }
      });
    }
    if (text.length) {
      editor.replaceRange(text.join('\n'), startPos, endPos);
    } else {
      editor.replaceRange('', startPos, endPos);
    }
    editor.setReadOnly(false);
    editor.focus();
  };
  const paste = async (event) => {
    const clipboard = event.clipboardData;

    const bool = verifyVideoSize(clipboard.files);

    if (bool) {
      event.preventDefault();
      const startPos = editor.getCursor();
      const endPos = { ...startPos, ch: startPos.ch + loadingText.length };

      editor.replaceSelection(loadingText);
      editor.setReadOnly(true);
      const urls = await upload(clipboard.files);
      const text = urls.map(({ name, url }) => {
        return `[${name}](${url})`;
      });

      editor.replaceRange(text.join('\n'), startPos, endPos);
      editor.setReadOnly(false);
      editor.focus();
    }
  };

  useEffect(() => {
    editor?.on('dragenter', dragenter);
    editor?.on('dragover', dragover);
    editor?.on('drop', drop);
    editor?.on('paste', paste);
    return () => {
      editor?.off('dragenter', dragenter);
      editor?.off('dragover', dragover);
      editor?.off('drop', drop);
      editor?.off('paste', paste);
    };
  }, [editor]);

  useEffect(() => {
    if (link.value && link.type === 'drop') {
      handleClick();
    }
  }, [link.value]);

  const addLink = (ctx) => {
    context = ctx;
    setEditor(context.editor);
    const text = context.editor?.getSelection();

    setVideoName({ ...videoName, value: text });

    setVisible(true);
  };

  const onUpload = async (e) => {
    if (!editor) {
      return;
    }
    const files = e.target?.files || [];
    const bool = verifyVideoSize(files);

    if (!bool) {
      return;
    }

    uploadImage({ file: e.target.files[0], type: 'post' }).then((url) => {
      setLink({ ...link, value: url });
    });
  };

  const onHide = () => setVisible(false);
  const onExited = () => editor?.focus();

  const handleSelect = (tab) => {
    setCurrentTab(tab);
  };

  return (
    <ToolItem {...item} onClick={addLink}>
      <Modal
        show={visible}
        onHide={onHide}
        onExited={onExited}
        fullscreen="sm-down">
        <Modal.Header closeButton>
          <h5 className="mb-0">{t('video.add_video')}</h5>
        </Modal.Header>
        <Modal.Body>
          <Tabs onSelect={handleSelect}>
            {/* <Tab eventKey="localVideo" title={t('video.tab_video')}>
              <Form className="mt-3" onSubmit={handleClick}>
                <Form.Group controlId="editor.videoLink" className="mb-3">
                  <Form.Label>
                    {t('video.form_video.fields.file.label')}
                  </Form.Label>
                  <Form.Control
                    type="file"
                    onChange={onUpload}
                    isInvalid={currentTab === 'localVideo' && link.isInvalid}
                  />

                  <Form.Control.Feedback type="invalid">
                    {t('video.form_video.fields.file.msg.empty')}
                  </Form.Control.Feedback>
                </Form.Group>

                <Form.Group
                  controlId="editor.videoDescription"
                  className="mb-3">
                  <Form.Label>
                    {`${t('video.form_video.fields.desc.label')} ${t(
                      'optional',
                      {
                        keyPrefix: 'form',
                      },
                    )}`}
                  </Form.Label>
                  <Form.Control
                    type="text"
                    value={videoName.value}
                    onChange={(e) =>
                      setVideoName({ ...videoName, value: e.target.value })
                    }
                    isInvalid={videoName.isInvalid}
                  />
                </Form.Group>
              </Form>
            </Tab> */}

            <Tab eventKey="remoteVideo" title={t('video.tab_url')}>
              <Form className="mt-3" onSubmit={handleClick}>
                <Form.Group controlId="editor.videoUrl" className="mb-3">
                  <Form.Label>
                    {t('video.form_url.fields.url.label')}
                  </Form.Label>
                  <Form.Control
                    type="text"
                    value={link.value}
                    onChange={(e) =>
                      setLink({ ...link, value: e.target.value })
                    }
                    isInvalid={currentTab === 'remoteVideo' && link.isInvalid}
                  />
                  <Form.Control.Feedback type="invalid">
                    {t('video.form_url.fields.url.msg.empty')}
                  </Form.Control.Feedback>
                </Form.Group>

                <Form.Group controlId="editor.videoName" className="mb-3">
                  <Form.Label>
                    {`${t('video.form_url.fields.name.label')} ${t('optional', {
                      keyPrefix: 'form',
                    })}`}
                  </Form.Label>
                  <Form.Control
                    type="text"
                    value={videoName.value}
                    onChange={(e) =>
                      setVideoName({ ...videoName, value: e.target.value })
                    }
                    isInvalid={videoName.isInvalid}
                  />
                </Form.Group>
              </Form>
            </Tab>
          </Tabs>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="link" onClick={() => setVisible(false)}>
            {t('video.btn_cancel')}
          </Button>
          <Button variant="primary" onClick={handleClick}>
            {t('video.btn_confirm')}
          </Button>
        </Modal.Footer>
      </Modal>
    </ToolItem>
  );
};

export default memo(Video);
