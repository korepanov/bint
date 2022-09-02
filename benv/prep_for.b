#import "stdlib/core.b"

string root_source;

void init(){
	get_root_source(root_source);
	SET_SOURCE(root_source);
	SET_DEST("benv/prep_for_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
};

bool is_for(string command){
	stack s;
	string buf;

	buf = "for(";
	s = ops(command, buf);
	s.pop(buf);
	
	return (NOT("end" == buf));
};

void main(){
	string buf;
	int pos;
	int command_len;
	stack s;

	init();
	#next:
	switch_command();
	if (NOT("end" == command)){
		if (is_for(command)){
			send_command(command);
			switch_command();
			send_command(command);
			switch_command(); 
			send_command(command);
			switch_command();
			buf = "{";
			s = ops(command, buf);
			s.pop(buf);
			if ("end" == buf){
				println("prep_for: ERROR");
				exit(1);			
			};
			pos = int(buf); 
			s.pop(buf);
			if (NOT("end" == buf)){
				println("prep_for: ERROR");
				exit(1);			
			};
			pos = (pos + 1);
			buf = command[0:pos];
			buf = (buf + "print(\"\")");
			send_command(buf);
			command_len = len(command);
			buf = command[pos:command_len];
			send_command(buf);		
		}else{
			send_command(command);
		};
		goto(#next);	
	};
	finish();
};

main();
