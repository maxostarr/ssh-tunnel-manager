export namespace ssh_manager {
	
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
	export class SshManagerRemote {
	    id: string;
	    name: string;
	    host: string;
	    port: number;
	    username: string;
	    status: string;
	    tunnels: SshManagerTunnel[];
	
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
	        this.status = source["status"];
	        this.tunnels = this.convertValues(source["tunnels"], SshManagerTunnel);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SshManagerRemoteData {
	    id: string;
	    name: string;
	    host: string;
	    port: number;
	    username: string;
	    status: string;
	
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
	        this.status = source["status"];
	    }
	}

}

