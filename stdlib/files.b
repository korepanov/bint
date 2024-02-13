/*
 deletes file in file_path
*/
void del_file(string file_path){
	if (exists(file_path)){
		$del_f(file_path);	
	}else{
		error = ("no such file: " + file_path); 
		
		if ($toPanic){
			print((error + "\n"));
			exit(1); 		
		};	
	};
};
/* opens file with file_path and mode 
 
 modes:
 mode = read 
 mode = write
 mode = append 

 return file descriptor number  
*/

int open_file(string file_path, string mode){
	int res; 
	int m;

	if (mode == "read"){
		m = 0;
		if (NOT(exists(file_path))){
			error = ("could not open file: no such file " + file_path);
			if ($toPanic){
				print((error + "\n"));
				exit(1);			
			};			
			return -1;		
		};
	}else if (mode == "write"){
		m = 1;	

		if (exists(file_path)){
			del_file(file_path); 		
		};
	}else if (mode == "append"){
		m = 2;	
	}else{
		error = "no such mode for opening file";
		if ($toPanic){
			print((error + "\n"));
			exit(1);		
		};
		return res; 
	};

	res = $open_f(file_path, m); 
	return res;
};

/*
 reads size bytes from file with descriptor_number

 returns read info as a string 
*/
string read_file(int descriptor_number, int size){
	int bytes; 
	string s; 

	bytes = $read_f(descriptor_number, size, s);
	return s;
};

/*
 writes s string to file with descriptor_number
*/
void write_file(int descriptor_number, string s){
	int bytes;
	bytes = $write_f(descriptor_number, s);

	if (bytes < 0){
		error = ("could not write to file with file descriptor number " + str(descriptor_number)); 

		if ($toPanic){
			print((error + "\n"));
			exit(1);
		};
	};
};

/*
 closes file with descriptor_number 
*/
void close_file(int descriptor_number){
	bool $toPanicOld;
	$toPanicOld = $toPanic;
	
	$close_f(descriptor_number); 
	
	try{
		string s;
		int bytes; 
		int size;
		size = 1;
		bytes = $read_f(descriptor_number, size, s); 
	};

	
	if (error == ""){
		error = ("could not close file with file descriptor number " + str(descriptor_number));	
		
		if ($toPanicOld){
			print((error + "\n"));
			exit(1);		
		};	
	};
	
	error = "";
};
