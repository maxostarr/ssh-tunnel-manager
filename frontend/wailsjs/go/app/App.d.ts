// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {ssh_manager} from '../models';

export function AddRemote(arg1:string,arg2:string,arg3:number,arg4:string):Promise<boolean>;

export function AddTunnel(arg1:string,arg2:number,arg3:string,arg4:number):Promise<boolean>;

export function Connect(arg1:string):Promise<boolean>;

export function DeleteRemote(arg1:string):Promise<boolean>;

export function Disconnect(arg1:string):Promise<void>;

export function GetRemote(arg1:string):Promise<ssh_manager.SshManagerRemote>;

export function GetRemotes():Promise<Array<ssh_manager.SshManagerRemoteData>>;

export function GetTunnels(arg1:string):Promise<Array<ssh_manager.SshManagerTunnel>>;

export function RemoveTunnel(arg1:string,arg2:number):Promise<boolean>;

export function UpdateRemote(arg1:string,arg2:string,arg3:string,arg4:number,arg5:string):Promise<boolean>;

export function WithEvents():Promise<void>;
