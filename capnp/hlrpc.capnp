using Go = import "/go.capnp";
@0xb6de64c471ee8bfb;
$Go.package("hlrpc");
$Go.import("hlinspect/internal/hlrpc");

interface HalfLife {
    getFullPlayerState @0 () -> (state :FullPlayerState);
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
