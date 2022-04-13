import React, {ChangeEvent, useEffect, useState} from "react";
import Cookies from "universal-cookie";
import axios from "axios";

interface Image {
  id: number;
  name: string;
  size: number;
  description: string;
  s3Link: string;
  dateUploaded: Date;
  label: string;
}

const ImagePage: React.FC = () => {

  const [images, setImages] = useState<Image[]>();
  const [isLoad, setLoad] = useState<boolean>(true);

  const getImages = () => {
    setLoad(true);
    axios.get<Image[]>('http://localhost:8080/api/images', {
      headers: {
        'Authorization': `Bearer ${new Cookies().get('jwt')}`
      },
      params: {page: 1, pageSize: 10}
    })
      .then(resp => {
        setLoad(false);
        setImages(resp.data)
        console.log(resp)
      })
  }

  const onUploadFilepath = (event: ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files;
    if (files === null || files.length === 0) {
      return;
    }
    const config = {
      headers: {
        'content-type': 'multipart/form-data',
        'Authorization': `Bearer ${new Cookies().get('jwt')}`
      }
    }
    const formData = new FormData();
    formData.append('img', files[0]);
    axios.post('http://localhost:8080/api/images', formData, config)
      .then(v => {
        getImages();
        event.preventDefault();
      })
  }

  useEffect(getImages, [])

  return <>
    <div className="container">
      <form action="#">
        <div className="file-field input-field center">
          <div className="btn" >
            <i className="material-icons left">cloud_upload</i>Загрузить
            <input type="file" accept="image/*" onChange={(t) => onUploadFilepath(t)}/>
          </div>
          <div className="file-path-wrapper">
            <input className="file-path validate" type="text" onChange={() => console.log('hi2')}/>
          </div>
        </div>
      </form>
      <div className="row">
      {isLoad ? <h1 className="center">Loading...</h1> : images?.map(image => (
          <div className="col s4">
            <div className="card">
              <div className="card-image">
                <img height="200" width="auto" src={image.s3Link}/>
              </div>
              <div className="card-content">
                {image.label}
              </div>
              <div className="card-action">
                {image.name}
              </div>
              <div className="card-action">
                {image.dateUploaded}
              </div>
            </div>
          </div>
      ))}
      </div>
    </div>
  </>
}

export default ImagePage