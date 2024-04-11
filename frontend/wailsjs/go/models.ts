export namespace ssh_manager {
	
	
	export class SshManagerRemote {
	    id: string;
	    name: string;
	    host: string;
	    port: number;
	    username: string;
	
	    static createFrom(source: any = {}) {
	        return new SshManagerRemote(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.username = source["username"];
	    }
	}
	export class SshManagerRemoteData {
	    id: string;
	    name: string;
	    host: string;
	    port: number;
	    username: string;
	
	    static createFrom(source: any = {}) {
	        return new SshManagerRemoteData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.username = source["username"];
	    }
	}
	export class SshManagerTunnel {
	    id: string;
	    local_port: number;
	    remote_host: string;
	    remote_port: number;
	    remote_id: string;
	
	    static createFrom(source: any = {}) {
	        return new SshManagerTunnel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.local_port = source["local_port"];
	        this.remote_host = source["remote_host"];
	        this.remote_port = source["remote_port"];
	        this.remote_id = source["remote_id"];
	    }
	}

}

