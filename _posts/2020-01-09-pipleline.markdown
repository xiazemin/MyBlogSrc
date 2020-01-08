---
title: 流水线冒险及解决方法
layout: post
category: linux
author: 夏泽民
---
由于一段机器语言程序的邻近指令之间出现了某种关联后，为了避免出错而使得它们不能同时被解释的现象，又称相关冲突。

    在流水解释过程中可能会出现三种相关，这三种相关是资源相关、数据相关和控制相关。

    1. 资源相关

    资源相关----是指多条指令进入流水线后在同一个时钟周期内争用同一功能部件所发生的相关。

    在图1所示的流水解释时空图中，在第4个时钟周期时，第I1条指令的MEM段与第I4条指令的IF段都要访问存储器。当数据和指令混存在同一个存储器且只有一个访问端口时，便会发生两条指令争用同一个存储器资源的相关冲突。
<!-- more -->
 解决资源相关的方法主要有以下五种：

    (1)从时间上推后下一条指令的访存操作。

    (2)让操作数和指令分别存放于两个独立编址且可同时访问的主存储器中。

    (3)仍然维持指令和操作数混存，但采用多模块交叉主存结构。

    (4)在CPU内增设指令Cache。

    (5)在CPU内增设指令Cache和数据Cache。

    2. 数据相关

    数据相关----是指由于相邻的两条或多条指令使用了相同的数据地址而发生的关联。这里所说的数据地址包括存储单元地址和寄存器地址。

    例如，有如下三条指令依次流入如图2所示的流水线：

        ADD R1，R2，R3    ；(R2)+(R3)→R1

        SUB R1，R4，R5    ；(R4)-(R5)→R1

        AND R4，R1，R7    ；(R1)∧(R7)→R4
        
  第三条指令与前面两条指令之间存在着关于寄存器R1的先写后读（RAW）数据相关。

    解决数据相关的方法主要有以下两种：

    (1)推后相关单元的读。
      (2)设置相关专用通路，又称采用定向传送技术。
       推后相关单元的读和设置相关专用通路是解决流水解释方式中数据相关的两种基本方法。推后相关单元的读是以降低速度为代价，使设备基本上不增加，而设置相关专用通路是以增加硬件为代价，使流水解释的性能尽可能不降低。

    对流水的流动顺序的安排和控制可以有两种方式：顺序流动方式和异步流动方式。

    顺序流动方式----是指流水线输出端的任务（指令）流出顺序和输入端的流入顺序一样。

    顺序流动方式的优点是控制比较简单。其缺点是一旦发生数据相关后，在空间和时间上都会有损失，使得流水线的吞吐率和功能部件的利用率降低。

    异步流动方式----是指流水线输出端的任务（指令）流出顺序和输入端的流入顺序可以不一样。

    异步流动方式的优点是流水线的吞吐率和功能部件的利用率都不会下降。但采用异步流动方式带来了新的问题，采用异步流动的控制复杂，而且会发生在顺序流动中不会出现的其它相关。由于异步流动要改变指令的执行顺序，同时流水线的异步流动还会使相关情况复杂化，会出现除先写后读（RAW）相关以外的先读后写（WAR）相关和写写（WAW）相关。

        3. 控制相关

    控制相关----是指由转移指令引起的相关。

    解决控制相关的方法主要有以下两种：

    (1)延迟转移技术。

    延迟转移技术----将转移指令与其前面的与转移指令无关的一条或几条指令对换位置，让成功转移总是在紧跟的指令被执行之后发生，从而使预取的指令不作废。

    (2)转移预测技术。

    转移预测技术可分为静态转移预测和动态转移预测两种（由硬件来实现）。

    静态转移预测技术有两种实现方法：

    一种是分析程序结构本身的特点来进行预测。不少条件转移指令转向两个目标地址的概率是能够预估的，如若x86汇编语言中连续的两条语句为“ADD AX,BX”、“JNZ L1”，则转向标号L1的概率要高，只要编译程序把出现概率高的分支安排为猜选分支，就能显著减少由于处理条件转移所带来的流水线吞吐率和效率的损失。

    另一种是按照分支的方向来预测分支是否转移成功。一般来说，向后分支被假定为循环的一部分而且被假定为发生转移，这种静态预测的准确性是相当高的。向前分支被假定为条件语句的一部分而且被假定为不发生转移，这种静态预测的准确性比向后分支的准确性要低得多。

    要提高预测的准确度，可以采用动态预测的方法，在硬件上建立分支预测缓冲站及分支目标缓冲站，根据执行过程中转移的历史记录来动态的预测转移目标，其预测准确度可以提高到90%以上。这种方法已在现代微处理器的转移预测中得到了广泛应用。例如，Pentium 4微处理器中的L1 BTB（Branch Target Buffer，转移目标缓冲器）采用的就是动态转移预测技术，当转移指令不在BTB中时，则采用静态转移预测技术。

    在安腾体系结构设计时，采用了一种新的转移预测技术。该技术将传统的分支结构变为无分支的并行代码，当处理机在运行中遇到分支时，它并不是进行传统的分支预测并选择可能性最大的一个分支执行，而是利用多个功能部件按分支的所有可能的后续路径开始并行执行多段代码，并暂存各段代码的执行结果，直到处理机能够确认分支转移与否的条件是真是假时，处理机再把应该选择的路径上的指令执行结果保留下来。采用了这种技术后，可消除大部分转移指令对流水解释的影响，使得整个系统的运行速度得到提高。
    
    流水线冒险
回顾一下常用五阶段流水线：F（取指）D（译码）E（执行）M（访存）W（写回）
注意：对寄存器文件的写只有在时钟上升的时候才会更新！

数据冒险的原因：对寄存器文件的读写是在不同阶段进行的

1.用暂停来避免数据冒险
暂停时，处理器会停止流水线中一条或多条指令，直到冒险条件不再满足。
在本该正常处理的某指令阶段中，我们每次要把一条指令阻塞在译码阶段，就在执行阶段插入一个气泡。
实现简单，但性能不好，严重降低了整体的吞吐量
2.用转发来避免数据冒险
①译码阶段逻辑发现寄存器是操作数的源寄存器，但是在写端口上还有一个对该寄存器未进行的写，将结果值直接从一个流水线阶段传到较早阶段的技术称为数据转发（旁路）。这种情况的信号是W_valE。
②访存阶段中有对寄存器未进行的写，也可以用数据转发。这种情况的信号是m_valE。
③执行阶段ALU正在计算的值稍后会写入寄存器时，也可以将ALU的输出信号作为操作数valB，这种情况的信号是e_valM。
④可转发刚从内存中读出的值，这种情况的信号是m_valM。
⑤可转发写回阶段对端口M未进行的写，这种情况的信号是W_valM。
3.加载/使用数据冒险
有一类数据冒险不能单纯用转发来解决，因为内存读在流水线发生的比较晚，这时我们可以通过将暂停和转发结合起来，避免加载/使用数据冒险。
这种用暂停来处理加载/使用冒险的方法称为加载互锁，加载互锁和转发技术结合起来，足以处理所有可能类型的数据冒险
4.避免控制冒险
当处理器无法根据处于取值阶段的当前指令来确定下一条指令的地址时，就会出现控制冒险（在流水线化处理器中，控制冒险只会发生在ret指令和跳转指令）跳转指令只有在条件跳转方向预测错误时才会造成麻烦。
ret指令经过译码，执行和访存阶段时，流水线不能做任何有用的活动，我们只能在流水线中插入三个气泡，一旦ret指令到达写回阶段，PC选择逻辑就会将程序计数器设为返回地址。
跳转指令的执行阶段才会知道是否应该跳转，但在这之前流水线预测会选择分支并取指、译码、执行，当跳转指令通过执行阶段后，流水线会向译码阶段和执行阶段插入气泡，取消两条目标指令，同时取出跳转指令后面的那条指令，
缺点：两个时钟周期的指令处理能力是被浪费的。

通过慎重考虑流水线的控制逻辑，控制冒险是可以被处理的，暂停和往流水线中插入气泡的技术可以动态调整流水线的流程。

1 主要实现代码

1.1 PC

主体为一个32位寄存器，带有stall信号关闭写使能信号

module PC(IF_Result,Clk,En,Clrn,IF_Addr,stall);

input [31:0]IF_Result;

input Clk,En,Clrn,stall;

output [31:0] IF_Addr;

wire En_S;

assign En_S=En&~stall;

D_FFEC32 pc(IF_Result,Clk,En_S,Clrn,IF_Addr);

Endmodule

1.2 INSTMEM（指令存储器）

主要参考书上代码（附每条指令具体含义）

module INSTMEM(Addr,Inst);//指令存储器

 input[31:0]Addr;

 output[31:0]Inst;

 wire[31:0]Rom[31:0];

assign Rom[5'h00]=32'b001000_00000_00001_0000_0000_0000_0010; //addi$1,$0,2--$1=2  2 0 2

 assignRom[5'h01]=32'b001000_00000_00010_0000_0000_0000_0100; //addi $2,$0,4--$2=4  4 0 4

 assignRom[5'h02]=32'b000000_00010_00001_00011_00000_100010; //sub$3,$2,$1--$3=$2-$1=2  （$2,rs,mem级数据前推）2 4 2

 assignRom[5'h03]=32'b101011_00010_00001_0000_0000_0000_1010; //sw$1,10($2)--[10+$2]=$1=2 （$2,rs,wb级数据前推）e 4 a

 assign Rom[5'h04]=32'b100011_00010_00100_0000_0000_0000_1010;//lw $4,10($2)--$4=[10+$2]=2                 e4a

 assignRom[5'h05]=32'b000000_00100_00001_00011_00000_100000;//add $3,$4,$1--$3=$4+$1=6$3=6             422

 assignRom[5'h06]=32'b000000_00001_00011_00010_00000_100000;//add $2,$1,$3--$2=$1+$3=8 ($3 ,rt,mem) 624

 assignRom[5'h07]=32'b001000_00000_00001_0000_0000_0000_0010;//addi $1,$0,2--$1=2                                    202

 assignRom[5'h08]=32'b000000_00001_00010_00100_00000_100000;//add $4,$1,$2 --$4=$1+$2=10(&2 ,rt,wb)   8 2 6  assign Inst=Rom[Addr[6:2]];

endmodule

 

1.3 CONUNIT（控制部件）

module CONUNIT(M_Op,Op,Func,M_Z,Regrt,Se,Wreg,Aluqb,Aluc,Wmem,Pcsrc,Reg2reg,Rs,Rt,E_Rd,M_Rd,E_Wreg,M_Wreg,FwdA,FwdB,E_Reg2reg,stall);

input [5:0]Op,Func,M_Op;

input M_Z;

input E_Wreg,M_Wreg,E_Reg2reg;

input [4:0]E_Rd,M_Rd,Rs,Rt;

output Regrt,Se,Wreg,Aluqb,Wmem,Reg2reg,stall;

output [1:0]Pcsrc,Aluc;

output reg [1:0]FwdA,FwdB;

wire R_type=~|Op;

wireI_add=R_type&Func[5]&~Func[4]&~Func[3]&~Func[2]&~Func[1]&~Func[0];

wireI_sub=R_type&Func[5]&~Func[4]&~Func[3]&~Func[2]&Func[1]&~Func[0];

wireI_and=R_type&Func[5]&~Func[4]&~Func[3]&Func[2]&~Func[1]&~Func[0];

wireI_or=R_type&Func[5]&~Func[4]&~Func[3]&Func[2]&~Func[1]&Func[0];

wire I_addi=~Op[5]&~Op[4]&Op[3]&~Op[2]&~Op[1]&~Op[0];

wire I_andi=~Op[5]&~Op[4]&Op[3]&Op[2]&~Op[1]&~Op[0];

wire I_ori=~Op[5]&~Op[4]&Op[3]&Op[2]&~Op[1]&Op[0];

wire I_lw=Op[5]&~Op[4]&~Op[3]&~Op[2]&Op[1]&Op[0];

wire I_sw=Op[5]&~Op[4]&Op[3]&~Op[2]&Op[1]&Op[0];

wire I_beq=~Op[5]&~Op[4]&~Op[3]&Op[2]&~Op[1]&~Op[0];

wire I_bne=~Op[5]&~Op[4]&~Op[3]&Op[2]&~Op[1]&Op[0];

wireM_beq=~M_Op[5]&~M_Op[4]&~M_Op[3]&M_Op[2]&~M_Op[1]&~M_Op[0];

wireM_bne=~M_Op[5]&~M_Op[4]&~M_Op[3]&M_Op[2]&~M_Op[1]&M_Op[0];

wire I_J=~Op[5]&~Op[4]&~Op[3]&~Op[2]&Op[1]&~Op[0];

wire E_Inst = I_add|I_sub|I_and|I_or|I_sw|I_beq|I_bne;

assign Regrt = I_addi|I_andi|I_ori|I_lw|I_sw|I_beq|I_bne|I_J;

assign Se = I_addi|I_lw|I_sw|I_beq|I_bne;

assign Wreg = I_add|I_sub|I_and|I_or|I_addi|I_andi|I_ori|I_lw;

assign Aluqb = I_add|I_sub|I_and|I_or|I_beq|I_bne|I_J;

assign Aluc[1] = I_and|I_or|I_andi|I_ori;

assign Aluc[0] = I_sub|I_or|I_ori|I_beq|I_bne;

assign Wmem = I_sw;

assign Pcsrc[1] = (M_beq&M_Z)|(M_bne&~M_Z)|I_J;

assign Pcsrc[0] = I_J;

assign Reg2reg = I_add|I_sub|I_and|I_or|I_addi|I_andi|I_ori|I_sw|I_beq|I_bne|I_J;

always@(E_Rd,M_Rd,E_Wreg,M_Wreg,Rs,Rt)begin

    FwdA=2'b00;

   if((Rs==E_Rd)&(E_Rd!=0)&(E_Wreg==1))begin

        FwdA=2'b10;

    end else begin

       if((Rs==M_Rd)&(M_Rd!=0)&(M_Wreg==1))begin

            FwdA=2'b01;

        end

    end

end

always@(E_Rd,M_Rd,E_Wreg,M_Wreg,Rs,Rt)begin

    FwdB=2'b00;

   if((Rt==E_Rd)&(E_Rd!=0)&(E_Wreg==1))begin

               FwdB=2'b10;

    end else begin

       if((Rt==M_Rd)&(M_Rd!=0)&(M_Wreg==1))begin

            FwdB=2'b01;

        end

    end

end

assignstall=((Rs==E_Rd)|(Rt==E_Rd))&(E_Reg2reg==0)&(E_Rd!=0)&(E_Wreg==1);

endmodule


 

1.4 REGFILE(寄存器)

主要参考书上代码

module REGFILE(Ra,Rb,D,Wr,We,Clk,Clrn,Qa,Qb);
 
input [4:0]Ra,Rb,Wr;
 
input [31:0]D;
 
input We,Clk,Clrn;
 
output [31:0]Qa,Qb;
 
wire[31:0]Y_mux,Q31_reg32,Q30_reg32,Q29_reg32,Q28_reg32,Q27_reg32,Q26_reg32,Q25_reg32,Q24_reg32,Q23_reg32,Q22_reg32,Q21_reg32,Q20_reg32,Q19_reg32,Q18_reg32,Q17_reg32,Q16_reg32,Q15_reg32,Q14_reg32,Q13_reg32,Q12_reg32,Q11_reg32,Q10_reg32,Q9_reg32,Q8_reg32,Q7_reg32,Q6_reg32,Q5_reg32,Q4_reg32,Q3_reg32,Q2_reg32,Q1_reg32,Q0_reg32;
 
 
 
DEC5T32E dec(Wr,We,Y_mux);
 
 
 
REG32_DOWN A(D,Y_mux,Clk,Clrn,Q31_reg32,Q30_reg32,Q29_reg32,Q28_reg32,Q27_reg32,Q26_reg32,Q25_reg32,Q24_reg32,Q23_reg32,Q22_reg32,Q21_reg32,Q20_reg32,Q19_reg32,Q18_reg32,Q17_reg32,Q16_reg32,Q15_reg32,Q14_reg32,Q13_reg32,Q12_reg32,Q11_reg32,Q10_reg32,Q9_reg32,Q8_reg32,Q7_reg32,Q6_reg32,Q5_reg32,Q4_reg32,Q3_reg32,Q2_reg32,Q1_reg32,Q0_reg32);
 
 
 
MUX32X32 select1(Q0_reg32,Q1_reg32,Q2_reg32,Q3_reg32,Q4_reg32,Q5_reg32,Q6_reg32,Q7_reg32,Q8_reg32,Q9_reg32,Q10_reg32,Q11_reg32,Q12_reg32,Q13_reg32,Q14_reg32,Q15_reg32,Q16_reg32,Q17_reg32,Q18_reg32,Q19_reg32,Q20_reg32,Q21_reg32,Q22_reg32,Q23_reg32,Q24_reg32,Q25_reg32,Q26_reg32,Q27_reg32,Q28_reg32,Q29_reg32,Q30_reg32,Q31_reg32,Ra,Qa);
 
MUX32X32select2(Q0_reg32,Q1_reg32,Q2_reg32,Q3_reg32,Q4_reg32,Q5_reg32,Q6_reg32,Q7_reg32,Q8_reg32,Q9_reg32,Q10_reg32,Q11_reg32,Q12_reg32,Q13_reg32,Q14_reg32,Q15_reg32,Q16_reg32,Q17_reg32,Q18_reg32,Q19_reg32,Q20_reg32,Q21_reg32,Q22_reg32,Q23_reg32,Q24_reg32,Q25_reg32,Q26_reg32,Q27_reg32,Q28_reg32,Q29_reg32,Q30_reg32,Q31_reg32,Rb,Qb);
 
 
 
Endmodule


 

实现下降沿写入数据主要由将D_FF部件改为下降沿更新：

module D_FF_DOWN(D,Clk,Q,Qn);

input D,Clk;

output Q,Qn;

wire Clkn,Q0,Qn0;

not i0(Clkn,Clk);

D_Latch d0(D,Clk,Q0,Qn0);

D_Latch d1(Q0,Clkn,Q,Qn);

endmodule


1.5 ALU（计算部件）

仅实现加、减、按位与、按位或

module ALU(X,Y,Aluc,R,Z);//ALU代码

 input [31:0]X,Y;

 input [1:0]Aluc;

 output [31:0]R;

 output Z;

 wire[31:0]d_as,d_and,d_or,d_and_or;

 ADDSUB_32 as(X,Y,Aluc[0],d_as);

 assign d_and=X&Y;

 assign d_or=X|Y;

 MUX2X32select1(d_and,d_or,Aluc[0],d_and_or);

 MUX2X32seleted(d_as,d_and_or,Aluc[1],R);

assign Z=~|R;

endmodule

 


1.6 DATAMEM（数据存储器）

module DATAMEM(Addr,Din,Clk,We,Dout);
 
 input[31:0]Addr,Din;
 
 input Clk,We;
 
 output[31:0]Dout;
 
 reg[31:0]Ram[31:0];
 
 assign Dout=Ram[Addr[6:2]];
 
 always@(posedge Clk)begin
 
  if(We)Ram[Addr[6:2]]<=Din;
 
 end
 
 integer i;
 
 initial begin
 
  for(i=0;i<32;i=i+1)
 
   Ram[i]=0;
 
  end
 
endmodule


 

1.7 中间寄存器（以ID/EX级为例）

主要为若干个若干位寄存器组成，在加上stall信号控制清零。

moduleREGidex(D0,D1,D2,D3,D4,D5,D6,D7,D8,D9,D10,En,Clk,Clrn_S,Q0,Q1,Q2,Q3,Q4,Q5,Q6,Q7,Q8,Q9,Q10,D11,D12,Q11,Q12,stall);

 input [31:0] D6,D7,D8,D9;

 input [5:0]D3;

 input [4:0]D10;

 input [1:0]D4,D11,D12;

 input D0,D1,D2,D5;

 

 input En,Clk,Clrn_S,stall;

wire Clrn;

assign Clrn=Clrn_S&~stall;

 

 output [31:0] Q6,Q7,Q8,Q9;

 output [5:0] Q3;

 output [4:0]Q10;

 output [1:0]Q4,Q11,Q12;

 output Q0,Q1,Q2,Q5;

 

 D_FFEC q0(D0,Clk,En,Clrn,Q0);

 D_FFEC q1(D1,Clk,En,Clrn,Q1);

 D_FFEC q2(D2,Clk,En,Clrn,Q2);

 D_FFEC6 q3(D3,Clk,En,Clrn,Q3);

 D_FFEC2 q4(D4,Clk,En,Clrn,Q4);

 D_FFEC q5(D5,Clk,En,Clrn,Q5);

 D_FFEC32 q6(D6,Clk,En,Clrn,Q6);

 D_FFEC32 q7(D7,Clk,En,Clrn,Q7);

 D_FFEC32 q8(D8,Clk,En,Clrn,Q8);

 D_FFEC32 q9(D9,Clk,En,Clrn,Q9);

 D_FFEC5 q10(D10,Clk,En,Clrn,Q10);

 D_FFEC2 q11(D11,Clk,En,Clrn,Q11);

 D_FFEC2 q12(D12,Clk,En,Clrn,Q12);

endmodule


 

1.8 MAIN（组装部分）

module MAIN(Clk,En,Clrn,IF_ADDR,EX_X,EX_Y,EX_R);

 

input Clk,En,Clrn;

output[31:0] IF_ADDR,EX_R,EX_X,EX_Y;

 

wire [31:0] IF_Result,IF_Addr,IF_PCadd4,IF_Inst,D,ID_Qa,ID_Qb,ID_PCadd4,ID_Inst;

wire [31:0]E_PC,E_R1,E_R2,E_I,E_I_L2,Y,E_R,EX_PC,M_PC,M_R,M_S,Dout,W_D,W_C,ID_EXTIMM,Alu_X,E_NUM;

wire[5:0] M_Op,E_Op;

wire [4:0] ID_Wr,W_Wr,E_Rd,M_Rd;

wire [1:0]Aluc,Pcsrc,E_Aluc,FwdA,FwdB,E_FwdA,E_FwdB;

wire M_Z,Regrt,Se,Wreg,Aluqb,Reg2reg,Wmem;

wireE_Wreg,E_Reg2reg,E_Wmem,E_Aluqb,E_Z,Cout,M_Wreg,M_Reg2reg,M_Wmem,W_Wreg,W_Reg2reg,stall;

 

//IF

 

MUX4X32 mux4x32(IF_PCadd4,0,M_PC,0,Pcsrc,IF_Result);

PC pc(IF_Result,Clk,En,Clrn,IF_Addr,stall);

PCadd4 pcadd4(IF_Addr,IF_PCadd4);

INSTMEM instmem(IF_Addr,IF_Inst);

 

REG_ifid ifid(IF_PCadd4,IF_Inst,En,Clk,Clrn,ID_PCadd4,ID_Inst,stall);

 

//ID

CONUNITconunit(M_Op,ID_Inst[31:26],ID_Inst[5:0],M_Z,Regrt,Se,Wreg,Aluqb,Aluc,Wmem,Pcsrc,Reg2reg,ID_Inst[25:21],ID_Inst[20:16],E_Rd,M_Rd,E_Wreg,M_Wreg,FwdA,FwdB,E_Reg2reg,stall);

MUX2X5 mux2x5(ID_Inst[15:11],ID_Inst[20:16],Regrt,ID_Wr);

EXT16T32 ext16t32(ID_Inst[15:0],Se,ID_EXTIMM);

REGFILEregfile(ID_Inst[25:21],ID_Inst[20:16],D,W_Wr,W_Wreg,Clk,Clrn,ID_Qa,ID_Qb);

 

REGidex idex(Wreg,Reg2reg,Wmem,ID_Inst[31:26],Aluc,Aluqb,ID_PCadd4,ID_Qa,ID_Qb,ID_EXTIMM,ID_Wr,En,Clk,Clrn,E_Wreg,E_Reg2reg,E_Wmem,E_Op,E_Aluc,E_Aluqb,E_PC,E_R1,E_R2,E_I,E_Rd,FwdA,FwdB,E_FwdA,E_FwdB,stall);

 

//EX

SHIFTER32_L2 shifter2(E_I,E_I_L2);

MUX4X32 mux4x32_ex_1(E_R1,D,M_R,0,E_FwdA,Alu_X);

MUX4X32 mux4x32_ex_2(E_R2,D,M_R,0,E_FwdB,E_NUM);

MUX2X32 mux2x321(E_I,E_NUM,E_Aluqb,Y);

ALU alu(Alu_X,Y,E_Aluc,E_R,E_Z);

CLA_32 cla_32(E_PC,E_I_L2,0,EX_PC,Cout);

 

REGexmemexmem(E_Wreg,E_Reg2reg,E_Wmem,E_Op,EX_PC,E_Z,E_R,E_R2,E_Rd,En,Clk,Clrn,M_Wreg,M_Reg2reg,M_Wmem,M_Op,M_PC,M_Z,M_R,M_S,M_Rd);

 

//MEM

DATAMEM datamem(M_R,M_S,Clk,M_Wmem,Dout);

 

REGmemwbmemwb(M_Wreg,M_Reg2reg,M_R,Dout,M_Rd,En,Clk,Clrn,W_Wreg,W_Reg2reg,W_D,W_C,W_Wr);

 

//WB

MUX2X32 mux2x322(W_C,W_D,W_Reg2reg,D);

 

assign IF_ADDR=IF_Addr;

assign EX_R=E_R;

assign EX_X=Alu_X;

assign EX_Y=Y;

endmodule


2 结构图




图2-1 流水线CPU结构图




3 测试代码及仿真结果

3.1 测试代码

module TEST;

reg Clk;

reg En;

reg Clrn;

 

wire [31:0] IF_ADDR;

wire [31:0] EX_R;

wire [31:0] EX_X;

wire [31:0] EX_Y;

MAIN uut(

.Clk(Clk),

.En(En),

.Clrn(Clrn),

.IF_ADDR(IF_ADDR),

.EX_R(EX_R),

.EX_X(EX_X),

.EX_Y(EX_Y)

);

initial begin

Clk=0;Clrn=0;En=1;

#10;

Clk=1;Clrn=0;

#10;

Clrn=1;

Clk=0;

forever #20 Clk=~Clk;

end

 

endmodule

3.2 仿真结果





3.3 仿真结果分析

结合指令储存器，主要观察EX_X,EX_Y,EX_R三个变量（ALU计算部件的X，Y，R端）

指令地址

指令代码

含义

Qa

Qb

R

00000000

001000_00000_00001_0000_0000_0000_0010

addi $1,$0,2--$1=2

$0(0)

2(立即数)

$1(2)

00000004

001000_00000_00010_0000_0000_0000_0100

addi $2,$0,4--$2=4

$0(0)

4(立即数)

$2(4)

00000008

000000_00010_00001_00011_00000_100010

sub $3,$2,$1  （$2,rs,mem)

$2(4)

$1(2)

$3(2)

0000000C

101011_00010_00001_0000_0000_0000_1010

sw $1,10($2) （$2,rs,wb）

$2(4)

10

e

00000010

100011_00010_00100_0000_0000_0000_1010

lw $4,10($2)   $4=2

$2(4)

10

e

00000014

000000_00100_00001_00011_00000_100000

add $3,$4,$1 

$4(2)

$1(2)

$3(4)

00000018

000000_00001_00011_00010_00000_100000

add $2,$1,$3 ($3 ,rt,mem)

$1(2)

$3(4)

$2(6)

0000001C

001000_00000_00001_0000_0000_0000_0010

addi $1,$0,2--$1=2                                   

$0(0)

2(立即数)

$1(2)

00000020

000000_00001_00010_00100_00000_100000

add $4,$1,$2 (&2 ,rt,wb)

$1(2)

$2(6)

$4(8)
