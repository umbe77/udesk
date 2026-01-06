export namespace main {
	
	export class ListItem {
	    icon: string;
	    text: string;
	    description: string;
	    original: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new ListItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.icon = source["icon"];
	        this.text = source["text"];
	        this.description = source["description"];
	        this.original = source["original"];
	    }
	}

}

