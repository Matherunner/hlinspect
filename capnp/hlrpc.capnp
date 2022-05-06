using Go = import "/go.capnp";
@0xb6de64c471ee8bfb;
$Go.package("schema");
$Go.import("hlinspect/internal/hlrpc/schema");

const buttonForward :UInt16 = 1;
const buttonBack :UInt16 = 2;
const buttodMoveRight :UInt16 = 4;
const buttonMoveLeft :UInt16 = 8;
const buttonRight :UInt16 = 16;
const buttonLeft :UInt16 = 32;
const buttonDuck :UInt16 = 64;
const buttonJump :UInt16 = 128;
const buttonUse :UInt16 = 256;
const buttonAttack1 :UInt16 = 512;
const buttonAttack2 :UInt16 = 1024;
const buttonReload :UInt16 = 2048;

interface HalfLife {
    getFullPlayerState @0 () -> (state :FullPlayerState);
    startInputControl @1 () -> ();
    stopInputControl @2 () -> ();
    inputStep @3 (cmd :CommandInput) -> (state: FullPlayerState);
}

struct CommandInput {
    buttons @0 :UInt16;

    forwardspeed @1 :Int16;
    sidespeed @2 :Int16;
    upspeed @3 :Int16;

    yawspeed @4 :Float32;
    pitchspeed @5 :Float32;
}

struct FullPlayerState {
    positionX @0 :Float32;
    positionY @1 :Float32;
    positionZ @2 :Float32;

    velocityX @3 :Float32;
    velocityY @4 :Float32;
    velocityZ @5 :Float32;

    baseVelocityX @6 :Float32;
    baseVelocityY @7 :Float32;
    baseVelocityZ @8 :Float32;

    yaw @9 :Float32;
    pitch @10 :Float32;
    roll @11 :Float32;

    punchYaw @12 :Float32;
    punchPitch @13 :Float32;
    punchRoll @14 :Float32;

    entityFriction @15 :Float32;
    entityGravity @16 :Float32;

    onGround @17 :Bool;
    duckState @18 :DuckState;
    waterLevel @19 :UInt8;
}

enum DuckState {
    standing @0;
    ducking @1;
    ducked @2;
}
