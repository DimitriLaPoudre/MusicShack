export interface RequestUser {
	username: string;
	password: string;
	bestQuality: boolean;
}

export interface RequestAdmin {
	password: string;
}

export interface RequestAdminPassword {
	oldPassword: string;
	newPassword: string;
}

export interface RequestInstance {
	api: string;
	url: string;
}

export interface RequestFollow {
	api: string;
	id: string;
}

export interface RequestDownload {
	api: string;
	type: string;
	id: string;
	quality: string;
}
