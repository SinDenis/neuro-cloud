import React from "react";
import {Navigate} from "react-router-dom";
import Cookies from "universal-cookie";

const RequireAuth: ({children}: { children: any }) => React.ReactElement<any, string | React.JSXElementConstructor<any>> = ({ children }) => {
  const jwt: string = new Cookies().get('jwt');
  return jwt === null ? <Navigate to="/login" replace /> : children;
}

export default RequireAuth