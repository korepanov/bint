$filestdlib/files.b$;
$trace1$;
intopen_file(stringfile_path,stringmode){intres;
$trace2$;
intm;
$trace3$;
if(mode=="read"){m=0;
$trace4$;
}elseif(mode=="write"){m=1;
$trace5$;
}elseif(mode=="append"){m=2;
$trace6$;
}else{error="no such mode for opening file";
$trace7$;
if($toPanic){print((error+"\n"));
$trace8$;
exit(1);
$trace9$;
};
$trace10$;
returnres;
$trace11$;
};
$trace12$;
res=$open_f(file_path,m);
$trace13$;
returnres;
$trace14$;
};
$trace15$;
stringread_file(intdescriptor_number,intsize){intbytes;
$trace16$;
strings;
$trace17$;
bytes=$read_f(descriptor_number,size,s);
$trace18$;
returns;
$trace19$;
};
$trace20$;
voidclose_file(intdescriptor_number){bool$toPanicOld;
$trace21$;
$toPanicOld=$toPanic;
$trace22$;
$close_f(descriptor_number);
$trace23$;
try{strings;
$trace24$;
intbytes;
$trace25$;
bytes=read_file(descriptor_number,1,s);
$trace26$;
};
$trace27$;
string$s0;
$s0=str(descriptor_number);
if(error==""){error=("could not close file with file descriptor number "+$s0);
UNDEFINE($s0);
$trace28$;
if($toPanicOld){print((error+"\n"));
$trace29$;
exit(1);
$trace30$;
};
$trace31$;
};
$trace32$;
error="";
$trace33$;
};
$fileprogram.b$;
$trace1$;
voidmain(){intf;
$trace2$;
f=open_file("/home/slava/a1.txt","read");
$trace3$;
strings;
$trace4$;
s=read_file(f,10);
$trace5$;
close_file(f);
$trace6$;
print((s+"\n"));
$trace7$;
};
$trace8$;
main();
