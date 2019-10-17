// Objective-C API for talking to github.com/incognitochain/incognito-chain/wasm/gomobile Go package.
//   gobind -lang=objc github.com/incognitochain/incognito-chain/wasm/gomobile
//
// File is generated by gobind. Do not edit.

#ifndef __Gomobile_H__
#define __Gomobile_H__

@import Foundation;
#include "ref.h"
#include "Universe.objc.h"


FOUNDATION_EXPORT NSString* _Nonnull GomobileDeriveSerialNumber(NSString* _Nullable args, NSError* _Nullable* _Nullable error);

/**
 * GenerateBLSKeyPairFromSeed generates BLS key pair from seed
 */
FOUNDATION_EXPORT NSString* _Nonnull GomobileGenerateBLSKeyPairFromSeed(NSString* _Nullable args);

/**
 * args: seed
 */
FOUNDATION_EXPORT NSString* _Nonnull GomobileGenerateKeyFromSeed(NSString* _Nullable seedB64Encoded, NSError* _Nullable* _Nullable error);

FOUNDATION_EXPORT NSString* _Nonnull GomobileInitBurningRequestTx(NSString* _Nullable args, NSError* _Nullable* _Nullable error);

// skipped function InitParamCreatePrivacyTokenTx with unsupported parameter or return types


// skipped function InitParamCreatePrivacyTx with unsupported parameter or return types


FOUNDATION_EXPORT NSString* _Nonnull GomobileInitPrivacyTokenTx(NSString* _Nullable args, NSError* _Nullable* _Nullable error);

FOUNDATION_EXPORT NSString* _Nonnull GomobileInitPrivacyTx(NSString* _Nullable args, NSError* _Nullable* _Nullable error);

FOUNDATION_EXPORT NSString* _Nonnull GomobileInitWithdrawRewardTx(NSString* _Nullable args, NSError* _Nullable* _Nullable error);

FOUNDATION_EXPORT NSString* _Nonnull GomobileRandomScalars(NSString* _Nullable n, NSError* _Nullable* _Nullable error);

FOUNDATION_EXPORT NSString* _Nonnull GomobileScalarMultBase(NSString* _Nullable scalarB64Encode, NSError* _Nullable* _Nullable error);

FOUNDATION_EXPORT NSString* _Nonnull GomobileStaking(NSString* _Nullable args, NSError* _Nullable* _Nullable error);

#endif
