import React from "react";

const AudioPage: React.FC = () => {

  return <>
    <div className="container">
      <form action="#">
        <div className="file-field input-field">
          <div className="btn">
            <i className="material-icons left">cloud_upload</i> Загрузить
            <input type="file"/>
          </div>
          <div className="file-path-wrapper">
            <input className="file-path validate" type="text"/>
          </div>
        </div>
      </form>
    </div>
  </>
}

export default AudioPage