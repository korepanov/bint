$filestdlib/files.b$;
$trace1$;
bool$ex0;
$ex0=exists(file_path);
voiddel_file(stringfile_path){if($ex0){$del_f(file_path);
UNDEFINE($ex0);
$trace2$;
}else{error=("no such file: "+file_path);
$trace3$;
if($toPanic){print((error+"\n"));
$trace4$;
exit(1);
$trace5$;
};
$trace6$;
};
$trace7$;
};
$trace8$;
intopen_file(stringfile_path,stringmode){intres;
$trace9$;
intm;
$trace10$;
if(mode=="read"){m=0;
$trace11$;
bool$ex0;
$ex0=exists(file_path);
if(NOT($ex0)){error=("could not open file: no such file "+file_path);
UNDEFINE($ex0);
$trace12$;
if($toPanic){print((error+"\n"));
$trace13$;
exit(1);
$trace14$;
};
$trace15$;
return-1;
$trace16$;
};
$trace17$;
}elseif(mode=="write"){m=1;
$trace18$;
bool$ex0;
$ex0=exists(file_path);
if($ex0){del_file(file_path);
UNDEFINE($ex0);
$trace19$;
};
$trace20$;
}elseif(mode=="append"){m=2;
$trace21$;
}else{error="no such mode for opening file";
$trace22$;
if($toPanic){print((error+"\n"));
$trace23$;
exit(1);
$trace24$;
};
$trace25$;
returnres;
$trace26$;
};
$trace27$;
res=$open_f(file_path,m);
$trace28$;
returnres;
$trace29$;
};
$trace30$;
stringread_file(intdescriptor_number,intsize){intbytes;
$trace31$;
strings;
$trace32$;
bytes=$read_f(descriptor_number,size,s);
$trace33$;
returns;
$trace34$;
};
$trace35$;
voidwrite_file(intdescriptor_number,strings){intbytes;
$trace36$;
bytes=$write_f(descriptor_number,s);
$trace37$;
string$s0;
$s0=str(descriptor_number);
if(bytes<0){error=("could not write to file with file descriptor number "+$s0);
UNDEFINE($s0);
$trace38$;
if($toPanic){print((error+"\n"));
$trace39$;
exit(1);
$trace40$;
};
$trace41$;
};
$trace42$;
};
$trace43$;
voidclose_file(intdescriptor_number){bool$toPanicOld;
$trace44$;
$toPanicOld=$toPanic;
$trace45$;
$close_f(descriptor_number);
$trace46$;
try{strings;
$trace47$;
intbytes;
$trace48$;
intsize;
$trace49$;
size=1;
$trace50$;
bytes=$read_f(descriptor_number,size,s);
$trace51$;
};
$trace52$;
string$s0;
$s0=str(descriptor_number);
if(error==""){error=("could not close file with file descriptor number "+$s0);
UNDEFINE($s0);
$trace53$;
if($toPanicOld){print((error+"\n"));
$trace54$;
exit(1);
$trace55$;
};
$trace56$;
};
$trace57$;
error="";
$trace58$;
};
$fileprogram.b$;
$trace1$;
voidmain(){intf;
$trace2$;
f=open_file("test.txt","read");
$trace3$;
del_file("test.txt");
$trace4$;
close_file(f);
$trace5$;
};
$trace6$;
main();
