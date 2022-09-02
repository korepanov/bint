#import "stdlib/core.b"

string root_source;

void init(){
	get_root_source(root_source);
	SET_SOURCE(root_source);
	SET_DEST("benv/for_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
};

bool is_for(string command){
	stack s;
	string buf;
	int pos;

	buf = "for(";
	s = ops(command, buf);
	s.pop(buf);

	if (NOT("end" == buf)){
		pos = int(buf);
		if (NOT(0 == pos)){
			println("for: ERROR");
			exit(1);	
		};	
	};

	return (NOT("end" == buf));
};

void main(){
	init();
	int counter;
	int command_len;
	#next:
	switch_command();
	if (NOT("end" == command)){
		if (is_for(command)){
			counter = block_end();
			command_len = len(command);
			command = command[4:command_len];
			send_command(command);
			switch_command();
			send_command(command);
			switch_command();
			println(command); 
		};
		send_command(command);
		goto(#next);	
	};
	finish();
};

main();
