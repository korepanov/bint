#import "stdlib/core.b"

string root_source;

void init(){
	root_source = "benv/prep_if_program.b";
	SET_SOURCE(root_source);
	SET_DEST("benv/prep_try_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
	//DEL_DEST(root_source);
}; 


string modify_command(string command){ 
	stack st;
	string op; 
	op = "try{";
 	
	st = ops(command, op);
	string buf; 
	st.pop(buf); 
	if (buf == "end"){
		return command;	
	};

	int pos;
	pos = int(buf);
	
	buf = ((command[0:pos] + op) + "print(\"\")");
	send_command(buf);
	
	return command[(pos + 4):len(command)]; 
};

void main(){
	init();
	
	string command;
	next_command(command);

	while (NOT(command == "end")){
		send_command(modify_command(command)); 
		next_command(command);	
	};
	


	finish(); 
};

main();
